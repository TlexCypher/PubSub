package main

import (
	"net/http"
)

type Server interface {
	Run()
}

type ApplicationServer struct {
	srv *http.Server
}

func NewApplicationServer(srv *http.Server) *ApplicationServer {
	return &ApplicationServer{
		srv: srv,
	}
}

func (s *ApplicationServer) Run(configs map[string]func(http.ResponseWriter, *http.Request)) {
	for k, v := range configs {
		http.HandleFunc(k, v)
	}
	s.srv.ListenAndServe()
}
