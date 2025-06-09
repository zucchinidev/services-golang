package mid

import (
	"context"

	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

// Logger is a middleware that logs information about the request to the logs.
func Logger(ctx context.Context, log *logger.Logger, path, rawQuery, method, remoteAddr string, handler Handler) error {

	values := web.GetValues(ctx)

	if rawQuery != "" {
		path = path + "?" + rawQuery
	}

	// TraceID is already logged by our fundational layer.
	log.Info(ctx, "request started", "method", method, "path", path, "remoteAddr", remoteAddr)

	err := handler(ctx)

	log.Info(ctx, "request completed", "method", method, "path", path, "remoteAddr", remoteAddr, "statuscode", values.StatusCode)

	return err
}
