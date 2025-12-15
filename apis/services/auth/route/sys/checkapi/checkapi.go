// Package checkapi provides support to check the API health.
package checkapi

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/foundation/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
}
