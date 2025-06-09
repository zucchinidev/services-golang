package mid

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/app/mid"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

// Logger is a middleware that logs the request and response.
// This middleware belongs to the API layer since it's protocol-specific (HTTP)
func Logger(log *logger.Logger) web.MidHandler {
	mw := func(next web.Handler) web.Handler {
		// This is the middleware function that will be called for each request.
		// It will wrap the next handler and add logging functionality.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			handler := func(ctx context.Context) error {
				return next(ctx, w, r)
			}

			// We do not want to use protocol-specific code in the middleware in the app layer.
			// So we pass the handler to the app layer and let it handle the request.
			return mid.Logger(ctx, log, r.URL.Path, r.URL.RawQuery, r.Method, r.RemoteAddr, handler)
		}

		return h
	}

	return mw
}
