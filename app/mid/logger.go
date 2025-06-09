package mid

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

// Logger is a middleware that logs the request and response.
// This is a simple implementation so far
func Logger(log *logger.Logger) web.MidHandler {
	mw := func(next web.Handler) web.Handler {
		// This is the middleware function that will be called for each request.
		// It will wrap the next handler and add logging functionality.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			log.Info(ctx, "request started", "method", r.Method, "path", r.URL.Path)

			err := next(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", r.URL.Path)
			return err
		}

		return h
	}

	return mw
}
