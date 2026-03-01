package server

import (
	"net/http"
	"net/http/httptest"
	"recon/internal/repo"
	"strings"
	"testing"
)

func TestHandleListOK(t *testing.T) {
	repos := []*repo.Repo{
		{Name: "service-a", Path: "/tmp/service-a", Branch: "main"},
		{Name: "service-b", Path: "/tmp/service-b", Branch: "master"},
	}
	srv := New(repos)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.handleList(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandleListContainsRepoNames(t *testing.T) {
	repos := []*repo.Repo{
		{Name: "my-special-repo", Path: "/tmp/my-special-repo", Branch: "main"},
	}
	srv := New(repos)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.handleList(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "my-special-repo") {
		t.Errorf("expected response body to contain %q", "my-special-repo")
	}
}

func TestHandleListEmpty(t *testing.T) {
	srv := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.handleList(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 with empty repos, got %d", w.Code)
	}
}
