package mid

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/app/api/mid"
	"github.com/zucchini/services-golang/foundation/web"
)

// Metrics is a middleware that reports metrics to the application.
func Metrics() web.MidHandler {
	return func(next web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return mid.Metrics(ctx, func(ctx context.Context) error {
				return next(ctx, w, r)
			})
		}
	}
}
