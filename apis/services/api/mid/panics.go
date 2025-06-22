package mid

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/app/api/mid"
	"github.com/zucchini/services-golang/foundation/web"
)

// Panics is a middleware that recovers from panics and returns an error so it
// can be reported in Metrics and handled in Errors
func Panics() web.MidHandler {
	return func(next web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return mid.Panics(ctx, func(ctx context.Context) error {
				return next(ctx, w, r)
			})
		}
	}
}
