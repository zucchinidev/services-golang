package authapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/zucchini/services-golang/app/api/authclient"
	"github.com/zucchini/services-golang/app/api/errs"
	"github.com/zucchini/services-golang/app/api/mid"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/web"
)

type api struct {
	au *auth.Auth
}

func newAPI(au *auth.Auth) *api {
	return &api{au: au}
}

func (a *api) token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := web.Param(r, "kid")
	if kid == "" {
		return errs.New(errs.FailedPrecondition, errors.New("missing kid"))
	}

	claims := mid.GetClaims(ctx)

	tkn, err := a.au.GenerateToken(kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	// TODO we need to move this out of here since the protocol layer must call
	// the app layer. We hack this for now until we create the token package
	token := struct {
		Token string `json:"token"`
	}{
		Token: tkn,
	}

	return web.Respond(ctx, w, token, http.StatusOK)
}

func (a *api) authenticate(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	// This middleware is actually handling the authentication. So if the code
	// gets to this handler, authentication passed.

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := authclient.AuthenticateResp{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (a *api) authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var auth authclient.Authorize

	if err := web.Decode(r, &auth); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	if err := a.au.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims [%v], rule [%v] %v", auth.Claims, auth.Rule, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
