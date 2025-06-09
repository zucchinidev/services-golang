// Package mux provides support to bind domain level routes to handlers.
package mux

import (
	"os"

	"github.com/zucchini/services-golang/apis/services/sales/route/sys/checkapi"
	"github.com/zucchini/services-golang/foundation/web"
)

func WebAPI(shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown)

	checkapi.Routes(mux)

	return mux
}
