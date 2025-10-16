package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Router is wrapped chi router
type Router struct {
	chiRouter   *chi.Mux
	groupPrefix string
	middlewares []any
}

// NewRouter creates a chi router with some additional fields such as group prefix and additional middlewares.
func NewRouter() *Router {
	return &Router{chiRouter: chi.NewRouter()}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chiRouter.ServeHTTP(w, req)
}

// Group adds a prefix and middlewares for all routes registered inside registerRoutes function.
func (r *Router) Group(prefix string, registerRoutes func(), middlewares ...any) {
	previousGroupPrefix := r.groupPrefix
	previousMiddlewares := r.middlewares

	// Adding group prefix for all the routes, which we want to register
	if r.groupPrefix != "" {
		r.groupPrefix += "/" + strings.Trim(prefix, "/")
	} else {
		r.groupPrefix = strings.Trim(prefix, "/")
	}

	// Adding additional middlewares, which we want to use with the routes
	r.middlewares = append(r.middlewares, middlewares...)

	registerRoutes()

	// Return everything to the state it was before
	r.groupPrefix = previousGroupPrefix
	r.middlewares = previousMiddlewares
}

// convertToMiddleware converts any supported middleware type into func(http.Handler) http.Handler.
func convertToMiddleware(handlerObject any) func(http.Handler) http.Handler {
	if handlerFunc, ok := handlerObject.(func(next http.Handler) http.Handler); ok {
		return handlerFunc
	}

	if handlerFunc, ok := handlerObject.(func(next http.Handler) http.HandlerFunc); ok {
		return func(next http.Handler) http.Handler {
			return handlerFunc(next)
		}
	}

	if handlerFunc, ok := handlerObject.(func(http.ResponseWriter, *http.Request)); ok {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerFunc(w, r)
				next.ServeHTTP(w, r)
			})
		}
	}

	panic(fmt.Sprintf("unsupported middleware type: %T", handlerObject))
}

// wrapMiddlewares converts any middleware type into func(http.Handler) http.Handler.
func wrapMiddlewares(functions ...any) []func(http.Handler) http.Handler {
	middlewares := make([]func(http.Handler) http.Handler, 0, len(functions))
	for _, m := range functions {
		middlewares = append(middlewares, convertToMiddleware(m))
	}
	return middlewares
}

// Post registers a POST route with optional middleware.
func (r *Router) Post(path string, handler func(http.ResponseWriter, *http.Request), functions ...any) {
	fullPath := "/" + strings.Trim(r.groupPrefix+"/"+path, "/")
	allMiddlewares := append(r.middlewares, functions...)

	r.chiRouter.With(wrapMiddlewares(allMiddlewares...)...).Post(fullPath, handler)
}

// Delete registers a DELETE route with optional middleware.
func (r *Router) Delete(path string, handler func(http.ResponseWriter, *http.Request), functions ...any) {
	fullPath := "/" + strings.Trim(r.groupPrefix+"/"+path, "/")
	allMiddlewares := append(r.middlewares, functions...)

	r.chiRouter.With(wrapMiddlewares(allMiddlewares...)...).Delete(fullPath, handler)
}

// Get registers a GET route with optional middleware.
func (r *Router) Get(path string, handler func(http.ResponseWriter, *http.Request), functions ...any) {
	fullPath := "/" + strings.Trim(r.groupPrefix+"/"+path, "/")
	allMiddlewares := append(r.middlewares, functions...)

	r.chiRouter.With(wrapMiddlewares(allMiddlewares...)...).Get(fullPath, handler)
}
