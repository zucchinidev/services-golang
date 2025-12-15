package mux

import (
	"os"

	"github.com/zucchini/services-golang/apis/services/api/mid"
	"github.com/zucchini/services-golang/apis/services/auth/route/authapi"
	"github.com/zucchini/services-golang/apis/services/auth/route/sys/checkapi"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

// WebAPI construct an http.Handler will all application routes bound.
func WebAPI(log *logger.Logger, a *auth.Auth, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	checkapi.Routes(app)
	authapi.Routes(app, a)

	return app
}
