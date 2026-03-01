package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func makeRepo(t *testing.T, dir string) {
	t.Helper()
	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("PlainInit: %v", err)
	}
	w, err := r.Worktree()
	if err != nil {
		t.Fatalf("Worktree: %v", err)
	}
	f, err := os.Create(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	if _, err = w.Add("README.md"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	_, err = w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@test.com", When: time.Now()},
	})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
}

func TestScanFindsRepo(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "myrepo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}
	makeRepo(t, repoDir)

	s := NewScanner(nil)
	repos, err := s.Scan(root)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if repos[0].Name != "myrepo" {
		t.Errorf("expected name %q, got %q", "myrepo", repos[0].Name)
	}
}

func TestScanFindsMultipleRepos(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"service-a", "service-b", "service-c"} {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		makeRepo(t, dir)
	}

	s := NewScanner(nil)
	repos, err := s.Scan(root)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(repos) != 3 {
		t.Errorf("expected 3 repos, got %d", len(repos))
	}
}

func TestScanIgnoresDirectories(t *testing.T) {
	root := t.TempDir()

	repoDir := filepath.Join(root, "myrepo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}
	makeRepo(t, repoDir)

	ignoredDir := filepath.Join(root, "node_modules", "some-pkg")
	if err := os.MkdirAll(ignoredDir, 0755); err != nil {
		t.Fatal(err)
	}
	makeRepo(t, ignoredDir)

	s := NewScanner([]string{"node_modules"})
	repos, err := s.Scan(root)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("expected 1 repo (node_modules should be ignored), got %d", len(repos))
	}
}

func TestScanEmptyRoot(t *testing.T) {
	root := t.TempDir()

	s := NewScanner(nil)
	repos, err := s.Scan(root)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(repos) != 0 {
		t.Errorf("expected 0 repos, got %d", len(repos))
	}
}
