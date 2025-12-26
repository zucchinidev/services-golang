package mid

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/app/api/authclient"
	"github.com/zucchini/services-golang/app/api/mid"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/web"
)

func AuthenticateOnServer(a *authclient.Client) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.AuthenticateOnServer(ctx, a, r.Header.Get("authorization"), hdl)
		}

		return h
	}

	return m
}

func AuthenticateLocal(a *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.AuthenticateLocal(ctx, a, r.Header.Get("authorization"), hdl)
		}

		return h
	}

	return m
}
