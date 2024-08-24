package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/secnex/pwexchange/api/middlewares"
	"github.com/secnex/pwexchange/api/routes"
	"github.com/secnex/pwexchange/storage"
)

type Server struct {
	ID       string
	Name     string
	Host     string
	Port     int
	BasePath string
	Vault    *storage.Vault
	Proxy    *httputil.ReverseProxy
	Routes   []routes.Route
	MU       sync.Mutex
	Auth     *middlewares.Auth
}

func NewServer(id string, name string, host string, port int, basePath string, vault *storage.Vault, proxy *httputil.ReverseProxy, auth *middlewares.Auth) *Server {
	return &Server{
		ID:       id,
		Name:     name,
		Host:     host,
		Port:     port,
		BasePath: basePath,
		Vault:    vault,
		Proxy:    proxy,
		Auth:     auth,
	}
}

func (s *Server) AddRoute(route routes.Route) {
	s.MU.Lock()
	s.Routes = append(s.Routes, route)
	s.MU.Unlock()
}

func (s *Server) RunServer() {
	r := http.NewServeMux()

	s.MU.Lock()
	for _, route := range s.Routes {
		r.HandleFunc(fmt.Sprintf("%s/%s", s.BasePath, route.Path), route.Handler)
		log.Printf("New route registered: %s/%s\n", s.BasePath, route.Path)
	}

	s.MU.Unlock()

	loggedRouter := middlewares.LoggingMiddleware(r)
	authRouter := s.Auth.Authenticate(loggedRouter)
	log.Printf("Starting %s (%s) on %s:%d\n", s.Name, s.ID, s.Host, s.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), authRouter))
}
