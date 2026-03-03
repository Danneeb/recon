package server

import (
	"net/http"
	"recon/internal/repo"
)

type Server struct {
	repos []*repo.Repo
}

func New(repos []*repo.Repo) *Server {
	return &Server{repos: repos}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleList)
	mux.HandleFunc("/repo", s.handleDetailView)
	return http.ListenAndServe(addr, mux)
}
