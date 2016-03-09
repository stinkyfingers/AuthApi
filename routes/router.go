package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/stinkyfingers/AuthApi/middleware"

	"net/http"
)

type Router struct {
	*httprouter.Router
}

func (r *Router) HandleRoute(method string, pattern string, handler middleware.Handler) {
	r.Router.Handle(method, pattern, wrap(handler))
}

func New() *Router {
	r := &Router{httprouter.New()}
	r.Router.RedirectTrailingSlash = true
	for _, route := range routes {
		r.HandleRoute(route.Method, route.Pattern, route.Handler)
	}
	return r
}

func wrap(h middleware.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, r, ps)
	})
}

func wrapMiddleware(h http.Handler) middleware.Middleware {
	return middleware.Middleware{h}
}
