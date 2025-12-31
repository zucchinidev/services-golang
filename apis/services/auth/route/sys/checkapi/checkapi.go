// Package checkapi provides support to check the API health.
package checkapi

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zucchini/services-golang/business/sqldb"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

type api struct {
	log   *logger.Logger
	db    *sqlx.DB
	build string
}

func newAPI(build string, log *logger.Logger, db *sqlx.DB) *api {
	return &api{
		build: build,
		log:   log,
		db:    db,
	}
}

func (api *api) liveness(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := map[string]any{
		"status":     "up",
		"build":      api.build,
		"host":       host,
		"name":       os.Getenv("KUBERNETES_NAME"),
		"podIP":      os.Getenv("KUBERNETES_POD_IP"),
		"node":       os.Getenv("KUBERNETES_NODE_NAME"),
		"namespace":  os.Getenv("KUBERNETES_NAMESPACE"),
		"GOMAXPROCS": runtime.GOMAXPROCS(0),
	}
	return web.Respond(ctx, w, data, http.StatusOK)
}

func (api *api) readiness(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := sqldb.StatusCheck(ctx, api.db); err != nil {
		failureStatus := "db not ready"
		api.log.Info(ctx, "readiness failure", "status", failureStatus)
		return web.Respond(ctx, w, map[string]string{"status": failureStatus}, http.StatusInternalServerError)
	}

	return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
}
