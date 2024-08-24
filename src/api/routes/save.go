package routes

import "net/http"

type Route struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

func NewRoute(path string, handler func(http.ResponseWriter, *http.Request)) Route {
	return Route{
		Path:    path,
		Handler: handler,
	}
}
