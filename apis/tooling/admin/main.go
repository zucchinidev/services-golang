// Package admin is the tooling for the admin API.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	_ "embed"

	"github.com/golang-jwt/jwt/v5"
	"github.com/open-policy-agent/opa/v1/rego"
)

const (
	bits = 2048
)

var (
	// embedding the Open Policy Agent (OPA) rules for authentication
	//go:embed rego/authentication.rego
	opaAuthenticationPolicyRego string
)

func main() {

	// ------------------------------------------------------------------------------------------------
	// Generate the Private and Public Key if the subcommand is "genkey" or the JWT if the subcommand is "genjwt"
	// ------------------------------------------------------------------------------------------------

	if len(os.Args) < 2 {
		log.Fatalln("missing subcommand")
	}

	switch os.Args[1] {
	case "genkey":
		if err := GenKey(); err != nil {
			log.Fatalln(err)
		}
	case "genjwt":
		role := "ADMIN"
		if len(os.Args) == 3 {
			role = strings.ToUpper(os.Args[2])
		}

		if !slices.Contains([]string{"ADMIN", "USER"}, role) {
			log.Fatalln("invalid role , use ADMIN or USER")
		}

		fmt.Printf("Generating JWT for role: %s\n", role)
		if err := GenJWT(role); err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatalln("invalid subcommand")
	}
}

func GenJWT(role string) error {

	// ------------------------------------------------------------------------------------------------
	// Generate the JWT
	// ------------------------------------------------------------------------------------------------

	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT is not valid
	// iat (issued at): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier of the JWT; can be used to prevent the JWT from being replayed

	// Payload section
	// ------------------------------------------------------------------------------------------------

	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "4801b850-e70f-4b1f-8fa7-d98aa2dac6d1",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8760)), // 1 year
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Roles: []string{role},
	}

	method := jwt.SigningMethodRS256 // RSA Signature with SHA-256, algorithm RS256 used for signing the JWT
	token := jwt.NewWithClaims(method, claims)

	// Headers section:
	// The headers section is used to specify data related to the token itself not to the payload
	// We are going to store the kid (key ID) to identify the key that was used to sign the JWT
	// Sometimes we need to rotate the keys, so we can use the kid to identify the key that was used to sign the JWT
	// The public key is used to verify the signature of the JWT, so when we need to verify the JWT we need to use
	// the public key associated with the private key that was used to sign the JWT
	token.Header["kid"] = "dc75a316-e862-45ca-a48b-0d67f229d62b"

	// ------------------------------------------------------------------------------------------------
	// Sign the JWT with the private Key
	// We are going to execute this from the Root of the project, so we are safe to use the relative path
	// ------------------------------------------------------------------------------------------------
	privateKeyPEM, err := os.ReadFile("zarf/keys/dc75a316-e862-45ca-a48b-0d67f229d62b.pem")
	if err != nil {
		return fmt.Errorf("unable to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %w", err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign token: %w", err)
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", tokenString)

	// ------------------------------------------------------------------------------------------------
	// Only for testing and debugging purposes, we are going to print also the private key
	// ------------------------------------------------------------------------------------------------

	ans1Bytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return fmt.Errorf("unable to marshal public key: %w", err)
	}

	ans1Block := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: ans1Bytes,
	}

	ans1PEM := pem.EncodeToMemory(&ans1Block)

	fmt.Println(string(ans1PEM))

	if err := checkPolicy(ans1PEM, tokenString, "service project"); err != nil {
		return fmt.Errorf("unable to check policy: %w", err)
	}

	return nil
}

func GenKey() error {

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	// Construct a PEM block for the private key
	// It represents the PEM encoding data structure
	pemPrivateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the Private Key File using the PEM block
	if err := pem.Encode(privateKeyFile, &pemPrivateBlock); err != nil {
		return err
	}

	// ------------------------------------------------------------------------------------------------
	// Generate the Public Key
	// ------------------------------------------------------------------------------------------------

	publicKeyFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	// Marshal the public key from the private key to PKIX, ASN.1 DER form
	derBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	// Construct a PEM block for the public key
	pemPublicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}

	// Write the public key to the Public Key File using the PEM block
	if err := pem.Encode(publicKeyFile, &pemPublicBlock); err != nil {
		return err
	}

	fmt.Println("Private Key File: private.pem")
	fmt.Println("Public Key File: public.pem")

	return nil
}

func checkPolicy(publicKeyPEM []byte, token string, issuer string) error {

	// ------------------------------------------------------------------------------------------------
	// Generate the OPA Policy
	// ------------------------------------------------------------------------------------------------

	ctx := context.Background()

	// ------------------------------------------------------------------------------------------------
	// Prepare the OPA Policy
	// ------------------------------------------------------------------------------------------------

	fmt.Println("OPA Authentication Policy:")
	fmt.Println(opaAuthenticationPolicyRego)

	query, err := rego.New(
		rego.Query(`x = data.sales.rego.auth`),
		rego.Module("policy.rego", opaAuthenticationPolicyRego),
	).PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("unable to prepare for evaluation: %w", err)
	}

	// ------------------------------------------------------------------------------------------------
	// Prepare the input for the OPA Policy
	// ------------------------------------------------------------------------------------------------

	input := map[string]any{
		"Key":   string(publicKeyPEM),
		"Token": token,
		"ISS":   issuer,
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("unable to evaluate policy: %w", err)
	}

	if len(results) == 0 {
		return fmt.Errorf("no results")
	}

	// check the value of the binding x, that is to say, data.sales.rego.auth
	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%v] ok[%v]", result, ok)
	}

	fmt.Println("TOKEN VALIDATION: Policy check passed")

	return nil
}
