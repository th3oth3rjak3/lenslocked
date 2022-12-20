package routers

import (
	mw "lenslocked/middleware"

	"github.com/go-chi/chi/v5"
)

type AppMiddleware struct {
	RequireUser *mw.RequireUser
	UserMW      *mw.User
}

type AppRouter struct {
	Router     *chi.Mux
	Middleware AppMiddleware
}

type RouteFunction func(router *AppRouter)
