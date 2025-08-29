// Package mid contains the middleware for the application.
package mid

import (
	"context"

	"github.com/zucchini/services-golang/app/api/errs"
	"github.com/zucchini/services-golang/foundation/logger"
)

// Errors handles errors coming out of the call chain. It detects normal application errors
// which are used to respond to the client in a uniform way.
func Errors(ctx context.Context, log *logger.Logger, next Handler) error {
	err := next(ctx)

	if err == nil {
		return nil
	}

	log.Error(ctx, "message", "ERROR", err.Error())

	if errs.IsError(err) {
		return errs.GetError(err)
	}

	return errs.Newf(errs.Unknown, "UNEXPECTED ERROR: %s", errs.Unknown.String())
}
