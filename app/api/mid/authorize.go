package mid

import (
	"context"

	"github.com/zucchini/services-golang/app/api/authclient"
	"github.com/zucchini/services-golang/app/api/errs"
)

func AuthorizeOnService(ctx context.Context, a *authclient.Client, rule string, handler Handler) error {
	userID, err := GetUserID(ctx)
	if err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: %v", err)
	}

	claims := GetClaims(ctx)
	authorize := authclient.Authorize{
		Claims: claims,
		UserID: userID,
		Rule:   rule,
	}
	if err := a.Authorize(ctx, authorize); err != nil {
		return errs.Newf(
			errs.Unauthenticated,
			"authorize: you are not authorized for that action, claims[%v] userID[%v] rule[%v]: %v",
			authorize.Claims,
			authorize.UserID,
			authorize.Rule,
			err,
		)
	}

	return handler(ctx)
}
