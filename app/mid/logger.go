package mid

import (
	"context"

	"github.com/zucchini/services-golang/foundation/logger"
)

// Logger is a middleware that logs information about the request to the logs.
func Logger(ctx context.Context, log *logger.Logger, path, rawQuery, method, remoteAddr string, handler Handler) error {

	if rawQuery != "" {
		path = path + "?" + rawQuery
	}

	log.Info(ctx, "request started", "method", method, "path", path, "remoteAddr", remoteAddr)

	err := handler(ctx)

	log.Info(ctx, "request completed", "method", method, "path", path, "remoteAddr", remoteAddr)

	return err
}
