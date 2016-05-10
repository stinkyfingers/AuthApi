package router

import (
	"github.com/stinkyfingers/AuthApi/database"
	"github.com/stinkyfingers/AuthApi/handlers"
	"github.com/stinkyfingers/AuthApi/middleware"

	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler middleware.Handler
}

var tokenAuthMiddleware = []middleware.Middleware{
	wrapMiddleware(middleware.TokenValidation{}),
}

var routes = []Route{
	Route{
		"Status",
		"GET",
		"/",
		middleware.Handler{Handle: status},
	},
	Route{
		"Create",
		"POST",
		"/create",
		middleware.Handler{Handle: handlers.Create},
	},
	Route{
		"Login",
		"POST",
		"/login",
		middleware.Handler{Handle: handlers.Login},
	},
	Route{
		"Logout",
		"POST",
		"/logout",
		middleware.Handler{Handle: handlers.Logout},
	},
	Route{
		"Authorize",
		"POST",
		"/authorize",
		middleware.Handler{Handle: handlers.Authorize, Middlewares: tokenAuthMiddleware},
	},
}

//status checks db connection and returns status
func status(ctx *middleware.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	status := "UP"
	err := database.Init()
	if err != nil {
		status = "DB error"
	}
	return status, err
}
