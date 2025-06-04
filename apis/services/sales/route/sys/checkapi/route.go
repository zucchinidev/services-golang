package checkapi

import (
	"net/http"
)

// Routes is the function that binds the checkapi routes to the mux.
func Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /liveness", liveness)
	mux.HandleFunc("GET /readiness", readiness)
}
