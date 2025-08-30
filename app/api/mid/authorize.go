package mid

import (
	"context"

	"github.com/zucchini/services-golang/app/api/errs"
	"github.com/zucchini/services-golang/business/api/auth"
)

func Authorize(ctx context.Context, a *auth.Auth, rule string, handler Handler) error {
	userID, err := GetUserID(ctx)
	if err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: %v", err)
	}

	claims := GetClaims(ctx)
	if err := a.Authorize(ctx, claims, userID, rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] userID[%v] rule[%v]: %v", claims, userID, rule, err)
	}

	return handler(ctx)
}
