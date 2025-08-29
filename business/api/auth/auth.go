package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/zucchini/services-golang/foundation/logger"
)

var ErrForbidden = errors.New("attempted action is not allowed")

type Claims struct {
	jwt.RegisteredClaims
	Roles []string
}

// HasRole checks if the claim has the given role
func (c *Claims) HasRole(role string) bool {
	return slices.Contains(c.Roles, role)
}

// KeyLookup is an interface for looking up keys by their identifier
// The return could be a PEM encoded string or a JWK based key.
type KeyLookup interface {
	PrivateKey(kid string) (key string, err error)
	PublicKey(kid string) (key string, err error)
}

// Config represents the information required to
// initialize auth.
type Config struct {
	Log       *logger.Logger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth is used to authenticate clients. It can generage a
// JWT for a set of user claims and recreate the claims
// by parsing the token.
type Auth struct {
	keyLookup KeyLookup
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
}

func New(cfg Config) *Auth {

	a := Auth{
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
	}

	return &a
}

func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)

	// Headers section:
	// The headers section is used to specify data related to the token itself not to the payload
	// We are going to store the kid (key ID) to identify the key that was used to sign the JWT
	// Sometimes we need to rotate the keys, so we can use the kid to identify the key that was used to sign the JWT
	// The public key is used to verify the signature of the JWT, so when we need to verify the JWT we need to use
	// the public key associated with the private key that was used to sign the JWT
	token.Header["kid"] = kid

	// ------------------------------------------------------------------------------------------------
	// Sign the JWT with the private Key
	// ------------------------------------------------------------------------------------------------
	privateKeyPEM, err := a.keyLookup.PrivateKey(kid)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %w", err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return tokenString, nil
}

// Authenticate processes the token to validate the sender's token is valid.
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, fmt.Errorf("unexpected authoriation header format, expected 'Bearer <token>'")
	}

	token := parts[1]
	var claims Claims
	// ParseUnverified is used to parse the token without verifying the signature
	// we want to delegate to OPA to verify the signature
	parsedToken, _, err := a.parser.ParseUnverified(token, &claims)
	if err != nil {
		return Claims{}, fmt.Errorf("unable to parse token: %w", err)
	}

	kidRaw, exists := parsedToken.Header["kid"]
	if !exists {
		return Claims{}, fmt.Errorf("no kid found in token")
	}

	kid, ok := kidRaw.(string)
	if !ok {
		return Claims{}, fmt.Errorf("kid malformed")
	}

	pem, err := a.keyLookup.PublicKey(kid)
	if err != nil {
		return Claims{}, fmt.Errorf("unable to lookup public key: %w", err)
	}

	// OPA will verify the token's signature
	input := map[string]any{
		"Token": token,
		"Key":   pem,
		"ISS":   a.issuer,
	}

	fmt.Printf("input: %+v\n", input)

	if err := opaPolictyEvaluation(ctx, regoScriptAuthentication, RuleAuthenticate, input); err != nil {
		return Claims{}, fmt.Errorf("unable to evaluate authentication policy: %w", err)
	}

	return claims, nil
}

// Authorize attempts to authorize the user with the provided input roles, if
// none of the input roles are within the user's claims, we return an error
// otherwise the user is authorized.
func (a *Auth) Authorize(ctx context.Context, claims Claims, userID uuid.UUID, rule string) error {

	input := map[string]any{
		"Roles":   claims.Roles,
		"Subject": claims.Subject,
		"UserID":  userID.String(),
		"Rule":    rule,
	}

	if err := opaPolictyEvaluation(ctx, regoScriptAuthorization, rule, input); err != nil {
		return fmt.Errorf("unable to evaluate authorization policy: %w", err)
	}

	// We must call the database here to see if the user is enabled.

	return nil
}

func opaPolictyEvaluation(ctx context.Context, regoScript string, rule string, input map[string]any) error {

	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", regoScript),
	).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%+v] ok[%v]", results, ok)
	}

	return nil
}
