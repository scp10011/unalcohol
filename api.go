package unalcohol

import "net/http"

type Middleware = func(*http.Request, http.ResponseWriter) error

type Server struct {
	Mux                       *http.ServeMux
	Addr                      string
	middleware                []Middleware
	StatusBadRequest          Middleware
	StatusMethodNotAllowed    Middleware
	StatusInternalServerError Middleware
}

type API interface {
	GetPath() string
	GetPtr() API
}

type BaseAPI struct {
	Path string
}

func (a *BaseAPI) GetPath() string {
	return a.Path
}

func (a *BaseAPI) GetPtr() API {
	return nil
}

func New() *Server {
	return &Server{
		Mux:                       http.NewServeMux(),
		middleware:                make([]Middleware, 0),
		StatusBadRequest:          DefaultStatusBadRequest,
		StatusMethodNotAllowed:    DefaultStatusMethodNotAllowed,
		StatusInternalServerError: DefaultStatusInternalServerError,
	}
}

func (srv *Server) GetMux() *http.ServeMux {
	return srv.Mux
}

func (srv *Server) Middleware(r *http.Request, resp http.ResponseWriter) error {
	for _, middleware := range srv.middleware {
		if err := middleware(r, resp); err != nil {
			return err
		}
	}
	return nil
}

func (srv *Server) Start() error {
	return http.ListenAndServe(srv.Addr, srv.Mux)
}
