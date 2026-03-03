package scanner

import (
	"io/fs"
	"log"
	"path/filepath"
	"recon/internal/repo"
	"slices"

	"github.com/go-git/go-git/v5"
)

type Scanner struct {
	ignore []string
}

func NewScanner(ignore []string) *Scanner {
	return &Scanner{ignore: ignore}
}

func (s *Scanner) Scan(root string) ([]*repo.Repo, error) {
	var gitRepos []*repo.Repo

	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if slices.Contains(s.ignore, d.Name()) {
				return fs.SkipDir
			}
		}

		if d.IsDir() && d.Name() == ".git" {
			repoPath, err := filepath.Abs(filepath.Dir(path))

			if err != nil {
				log.Printf("skipping %s: %v", path, err)
				return fs.SkipDir
			}
			r, err := git.PlainOpen(repoPath)
			if err != nil {
				log.Printf("skipping %s: %v", repoPath, err)
				return fs.SkipDir
			}
			branch, err := r.Head()
			if err != nil {
				log.Printf("skipping %s: %v", repoPath, err)
				return fs.SkipDir
			}
			gitRepos = append(gitRepos, &repo.Repo{Name: filepath.Base(repoPath), Path: repoPath, Branch: branch.Name().Short()})
			return fs.SkipDir
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return gitRepos, nil
}
