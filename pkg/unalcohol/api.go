package unalcohol

import "net/http"

type Server struct {
	Mux  *http.ServeMux
	Addr string
}

type API struct {
	Path string
}

func (a *API) SetPath(path string) *API {
	a.Path = path
	return a
}

func (srv *Server) GetMux() *http.ServeMux {
	return srv.Mux
}

func (srv *Server) Start() error {
	return http.ListenAndServe(srv.Addr, srv.Mux)
}
