package mid

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/app/api/mid"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/web"
)

func Authorize(a *auth.Auth, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Authorize(ctx, a, rule, hdl)
		}

		return h
	}

	return m
}
