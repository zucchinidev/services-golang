package checkapi

import (
	"github.com/zucchini/services-golang/apis/services/api/mid"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/web"
)

// Routes is the function that binds the checkapi routes to the mux.
func Routes(mux *web.App, a *auth.Auth) {

	authsMw := []web.MidHandler{mid.Authenticate(a), mid.Authorize(a, auth.RuleAdminOnly)}

	mux.HandleFuncNoMiddleware("GET /liveness", liveness)
	mux.HandleFuncNoMiddleware("GET /readiness", readiness)
	mux.HandleFunc("GET /testerror", testErr)
	mux.HandleFunc("GET /testpanic", testPanic)
	mux.HandleFunc("GET /testauth", liveness, authsMw...)
}
