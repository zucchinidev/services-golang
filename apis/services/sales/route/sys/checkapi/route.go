package checkapi

import "github.com/zucchini/services-golang/foundation/web"

// Routes is the function that binds the checkapi routes to the mux.
func Routes(mux *web.App) {
	mux.HandleFunc("GET /liveness", liveness)
	mux.HandleFunc("GET /readiness", readiness)
	mux.HandleFunc("GET /testerror", testErr)
}
