package mid

import (
	"context"
	"encoding/base64"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zucchini/services-golang/app/api/authclient"

	"github.com/zucchini/services-golang/app/api/errs"
	"github.com/zucchini/services-golang/business/api/auth"
)

func AuthenticateOnServer(ctx context.Context, authClient *authclient.Client, authorization string, handler Handler) error {
	resp, err := authClient.Authenticate(ctx, authorization)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	ctx = setUserID(ctx, resp.UserID)
	ctx = setClaims(ctx, resp.Claims)

	return handler(ctx)
}

func AuthenticateLocal(ctx context.Context, a *auth.Auth, authorization string, handler Handler) error {
	var err error

	parts := strings.Split(authorization, " ")

	switch parts[0] {
	case "Bearer":
		ctx, err = processJWT(ctx, a, authorization)
		if err != nil {
			return err
		}

	case "Basic":
		ctx, err = processBasic(ctx, authorization)
	}

	if err != nil {
		return err
	}

	return handler(ctx)
}

func processJWT(ctx context.Context, a *auth.Auth, token string) (context.Context, error) {
	claims, err := a.Authenticate(ctx, token)
	if err != nil {
		return ctx, err
	}

	if claims.Subject == "" {
		return ctx, errs.Newf(errs.Unauthenticated, "authenticate: you are not authorized for that action, no claims found")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.Newf(errs.Unauthenticated, "authenticate: parsing subject: %v", err)
	}

	ctx = setUserID(ctx, userID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

func processBasic(ctx context.Context, authorization string) (context.Context, error) {
	email, _, ok := parseBasicAuth(authorization)
	if !ok {
		return ctx, errs.Newf(errs.Unauthenticated, "authenticate: invalid Basic auth header")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return ctx, errs.Newf(errs.Unauthenticated, "authenticate: parsing email: %v", err)
	}

	// copy from tooling: "4801b850-e70f-4b1f-8fa7-d98aa2dac6d1"
	// we'll get this from the database
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "4801b850-e70f-4b1f-8fa7-d98aa2dac6d1",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.Newf(errs.Unauthenticated, "authenticate: parsing subject: %v", err)
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

func parseBasicAuth(authorization string) (string, string, bool) {
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}
