package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandleType struct {
	Method string
	Handle http.HandlerFunc
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]HandleType
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]HandleType),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandlers(path string, method string, handle http.HandlerFunc) {
	s.Handlers[path] = HandleType{Method: method, Handle: handle}
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {

		if handler.Method == http.MethodPost {
			s.Router.Post(path, handler.Handle)
		}

		if handler.Method == http.MethodGet {
			s.Router.Get(path, handler.Handle)
		}

	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
