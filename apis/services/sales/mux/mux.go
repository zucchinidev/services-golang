// Package mux provides support to bind domain level routes to handlers.
package mux

import (
	"os"

	"github.com/zucchini/services-golang/apis/services/api/mid"
	"github.com/zucchini/services-golang/apis/services/sales/route/sys/checkapi"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

func WebAPI(shutdown chan os.Signal, log *logger.Logger) *web.App {
	mux := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	checkapi.Routes(mux)

	return mux
}
