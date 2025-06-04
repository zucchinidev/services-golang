// Package debug provides handlers for the debugging endpoints.
package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/arl/statsviz"
)

// This is not definetily APP layer code because this is going to be a very heavily protocol driven

func Mux() *http.ServeMux {
	mux := http.NewServeMux()

	// -------------------------------------------------------------------------
	// Register debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	statsviz.Register(mux)

	return mux
}
