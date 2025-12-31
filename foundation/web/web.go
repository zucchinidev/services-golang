// Package web provides a small HTTP framework for building web applications.
package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
)

// Handler is a function that can handle a http request within our small little own HTTP
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint to our application that configures context for HTTP handlers.
// It wraps the standard library HTTP server to provide idiomatic Go patterns like
// context as first parameter and error returns.
type App struct {
	// We embed http.ServeMux so we can use it as our router.
	// My App is not everything that http.ServeMux is.
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

// NewApp creates a new App value that contains the information for the HTTP server.
func NewApp(shutdown chan os.Signal, mw ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

// SignalShutdown is used to gracefully Shutdown the app when an integrity issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) HandleFunc(pattern string, handler Handler, mw ...MidHandler) {
	// local middleware first
	// This allows us to for example, add an authentication middleware only to this handler.
	handler = wrapMiddleware(mw, handler)
	// general middleware afters
	handler = wrapMiddleware(a.mw, handler)

	a.ServeMux.HandleFunc(pattern, a.generateHandlerFunc(handler))
}

// HandleFuncNoMiddleware is a helper function that handles a http request
// without any middleware.
func (a *App) HandleFuncNoMiddleware(pattern string, handler Handler) {
	a.ServeMux.HandleFunc(pattern, a.generateHandlerFunc(handler))
}

func (a *App) generateHandlerFunc(handler Handler) http.HandlerFunc {

	h := func(w http.ResponseWriter, r *http.Request) {

		// Put any code here that needs to happen before the handler is called.
		// Remember, this is a fundational layer. Highly portable.
		// Do not log here or execute code specific to a handler.

		v := &Values{
			TraceID: uuid.NewString(),
			Now:     time.Now(),
		}

		ctx := setValues(r.Context(), v)

		if err := handler(ctx, w, r); err != nil {

			// This error could happen when we send the Shutdown signal or we cannot write down to the pipe.
			// This is because the manage the errors down to the handler level.

			if validateError(err) {
				// We prefer to Shutdown the server gracefully
				// rather than have the server in an inconsistent state.
				// It is a tough call, but I think this is the best choice.
				// It is better to restart than to have corrumpted file systems, data, etc.
				a.SignalShutdown()
				return
			}
		}

		// Put any code here that needs to happen after the handler is called.
	}

	return h
}

// validateError validates the error for special conditions that do not
// warrant an actual Shutdown by the system.
func validateError(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	// These errors are explicitly ignored (function returns false) because they represent client-side disconnections rather than server-side problems. They are:
	// Routine network issues that happen frequently in production systems
	// Not indicative of application logic problems
	// Expected behaviors when clients disconnect unexpectedly
	// Handled automatically by the TCP/IP stack

	switch {
	case errors.Is(err, syscall.EPIPE):
		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):
		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true

}
