// Package checkapi provides support to check the API health.
package checkapi

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/zucchini/services-golang/app/api/errs"
	"github.com/zucchini/services-golang/foundation/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
}

func testErr(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return errs.Newf(errs.FailedPrecondition, "test error - this error is safe")
	}

	return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
}

func testPanic(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		panic("test panic - we are panicking")
	}

	return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
}
