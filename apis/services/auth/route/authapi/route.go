package authapi

import (
	"github.com/zucchini/services-golang/apis/services/api/mid"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/web"
)

// Routes is the function that binds the authapi routes to the mux.
func Routes(mux *web.App, a *auth.Auth) {

	api := newAPI(a)
	authenticateLocal := mid.AuthenticateLocal(a)
	mux.HandleFunc("GET /auth/token/{kid}", api.token, authenticateLocal)
	mux.HandleFunc("GET /auth/authenticate", api.authenticate, authenticateLocal)
	mux.HandleFunc("POST /auth/authorize", api.authorize)
}
