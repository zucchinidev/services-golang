// Package mux provides support to bind domain level routes to handlers.
package mux

import "net/http"

// Every package in this project tries to create a kind of firewall or boundary.
// We need to organize APIs that communicate with different parts of the system. That's what a package does.
// We are not orgnizing code, but organizing APIs. Those APIs are taking to another level of this project by
// putting them into another layer of code and eventually, we will have a vertical layer of code that is part of our domain.
// In our core we will have package that provide purpose, not package that contain.
// That package API that provides purpose, will perform whatever data transformation we need.
// Every package needs to have thir type system. The type system is going to represent the data coming in and the data leaving that API
// When it comes to this type system, we can use concrete types or interfaces.
// Ideally, we want to use concrete types. Real Data!!
// We want to know what it is, it's a user, it's a product, it's a sale...
// When we need to use functions based on how it behaves instead of what it is, we will use interfaces. This is called polymorphism.
// Polymorphism means that a piece of code changes its behaviour depending upon the type of data that is operating on.
// Runtime polymorphism is when we use interfaces to achieve polymorphism.
// Static polymorphism is when we use generics to achieve polymorphism.
// General rule: use interfaces when they are needed as import types but never as return types. We want to leave the decauple of our data to the caller, in other words,
// for the person that needs it.
// In this case, we want to return a pointer to a concrete type.
func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/sales", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	return mux
}
