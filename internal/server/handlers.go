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

func (s *Server) handleDetailView(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	var found *repo.Repo

	for _, r := range s.repos {
		if r.Path == path {
			found = r
			break
		}

	}

	if found == nil {
		http.NotFound(w, r)
		return
	}

	branch := r.URL.Query().Get("branch")
	data, err := repo.GetRepoDetail(found, path, branch)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "detailView.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
