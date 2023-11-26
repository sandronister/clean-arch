package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandleType struct {
	Path   string
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
	s.Handlers[path+"|"+method] = HandleType{Path: path, Method: method, Handle: handle}
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, handler := range s.Handlers {

		if handler.Method == http.MethodPost {
			s.Router.Post(handler.Path, handler.Handle)
		}

		if handler.Method == http.MethodGet {
			s.Router.Get(handler.Path, handler.Handle)
		}

	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
