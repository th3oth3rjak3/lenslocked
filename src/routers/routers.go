package routers

import (
	mw "lenslocked/middleware"

	"github.com/labstack/echo/v4"
)

type AppMiddleware struct {
	RequireUser *mw.RequireUser
	UserMW      *mw.User
}

type AppRouter struct {
	Router     *echo.Echo
	Middleware AppMiddleware
}

type RouteFunction func(router *AppRouter)
