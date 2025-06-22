package checkapi

import "github.com/zucchini/services-golang/foundation/web"

// Routes is the function that binds the checkapi routes to the mux.
func Routes(mux *web.App) {
	mux.HandleFuncNoMiddleware("GET /liveness", liveness)
	mux.HandleFuncNoMiddleware("GET /readiness", readiness)
	mux.HandleFunc("GET /testerror", testErr)
	mux.HandleFunc("GET /testpanic", testPanic)
}
