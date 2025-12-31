package checkapi

import (
	"github.com/jmoiron/sqlx"
	"github.com/zucchini/services-golang/apis/services/api/mid"
	"github.com/zucchini/services-golang/app/api/authclient"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

// Routes is the function that binds the healing routes to the mux.
func Routes(build string, log *logger.Logger, mux *web.App, db *sqlx.DB, a *authclient.Client) {

	authsMw := []web.MidHandler{mid.AuthenticateOnServer(a), mid.AuthorizeOnService(a, auth.RuleAdminOnly)}
	api := newAPI(build, log, db)

	mux.HandleFuncNoMiddleware("GET /liveness", api.liveness)
	mux.HandleFuncNoMiddleware("GET /readiness", api.readiness)
	mux.HandleFunc("GET /testerror", api.testErr)
	mux.HandleFunc("GET /testpanic", api.testPanic)
	mux.HandleFunc("GET /testauth", api.liveness, authsMw...)
}
