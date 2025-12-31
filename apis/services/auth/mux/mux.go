package mux

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/zucchini/services-golang/apis/services/api/mid"
	"github.com/zucchini/services-golang/apis/services/auth/route/authapi"
	"github.com/zucchini/services-golang/apis/services/auth/route/sys/checkapi"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

// WebAPI construct an http.Handler will all application routes bound.
func WebAPI(build string, log *logger.Logger, db *sqlx.DB, a *auth.Auth, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	checkapi.Routes(build, log, app, db)
	authapi.Routes(app, a)

	return app
}
