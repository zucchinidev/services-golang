package checkapi

import (
	"github.com/jmoiron/sqlx"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

// Routes is the function that binds the checkapi routes to the mux.
func Routes(build string, log *logger.Logger, mux *web.App, db *sqlx.DB) {
	api := newAPI(build, log, db)
	mux.HandleFuncNoMiddleware("GET /liveness", api.liveness)
	mux.HandleFuncNoMiddleware("GET /readiness", api.readiness)
}
