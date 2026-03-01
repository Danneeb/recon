package server

import (
	"html/template"
	"net/http"
	"recon/internal/repo"
	"recon/web"
)

var tmpl = template.Must(template.ParseFS(web.Templates, "templates/*.html"))

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Repos []*repo.Repo
	}{
		Repos: s.repos,
	}
	if err := tmpl.ExecuteTemplate(w, "list.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
