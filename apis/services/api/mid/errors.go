package mid

import (
	"context"
	"net/http"

	"github.com/zucchini/services-golang/app/api/errs"
	"github.com/zucchini/services-golang/app/api/mid"
	"github.com/zucchini/services-golang/foundation/logger"
	"github.com/zucchini/services-golang/foundation/web"
)

var codeStatus [17]int

// init maps out the error codes to http status codes.
func init() {
	codeStatus[errs.OK.Value()] = http.StatusOK
	codeStatus[errs.Canceled.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.Unknown.Value()] = http.StatusInternalServerError
	codeStatus[errs.InvalidArgument.Value()] = http.StatusBadRequest
	codeStatus[errs.DeadlineExceeded.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.NotFound.Value()] = http.StatusNotFound
	codeStatus[errs.AlreadyExists.Value()] = http.StatusConflict
	codeStatus[errs.PermissionDenied.Value()] = http.StatusForbidden
	codeStatus[errs.ResourceExhausted.Value()] = http.StatusTooManyRequests
	codeStatus[errs.FailedPrecondition.Value()] = http.StatusBadRequest
	codeStatus[errs.Aborted.Value()] = http.StatusConflict
	codeStatus[errs.OutOfRange.Value()] = http.StatusBadRequest
	codeStatus[errs.Unimplemented.Value()] = http.StatusNotImplemented
	codeStatus[errs.Internal.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unavailable.Value()] = http.StatusServiceUnavailable
	codeStatus[errs.DataLoss.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unauthenticated.Value()] = http.StatusUnauthorized
}

func Errors(log *logger.Logger) web.MidHandler {
	mw := func(next web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Encapsulate the handler to prevent HTTP protocol details from propagating
			// down through the application layers. This maintains clean separation of concerns.
			hdlr := func(ctx context.Context) error {
				return next(ctx, w, r)
			}

			// Application layer middleware. No protocol details!
			if err := mid.Errors(ctx, log, hdlr); err != nil {
				// We test this before going to production. We do not check ok.
				errs := err.(errs.Error)
				// Application layer code to protocol layer code
				code := codeStatus[errs.Code.Value()]
				if err = web.Respond(ctx, w, errs, code); err != nil {
					return err
				}

				// If the error is a shutdown error, we need to return it
				// back to the base of the handle to shut down the server.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}
	return mw
}
