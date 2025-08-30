package mid

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/zucchini/services-golang/app/api/errs"
	"github.com/zucchini/services-golang/business/api/auth"
)

func Authenticate(ctx context.Context, auth *auth.Auth, authorization string, handler Handler) error {

	parts := strings.Split(authorization, " ")

	switch parts[0] {
	case "Bearer":
		var err error
		ctx, err = processJWT(ctx, auth, authorization)
		if err != nil {
			return err
		}
	}

	return handler(ctx)
}

func processJWT(ctx context.Context, auth *auth.Auth, token string) (context.Context, error) {
	claims, err := auth.Authenticate(ctx, token)
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
