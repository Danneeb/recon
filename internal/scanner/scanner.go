package scanner

import (
	"io/fs"
	"path/filepath"
	"recon/internal/repo"

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

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			for _, ignored := range s.ignore {
				if d.Name() == ignored {
					return fs.SkipDir
				}
			}
		}

		if d.IsDir() && d.Name() == ".git" {
			repoPath := filepath.Dir(path)
			r, err := git.PlainOpen(repoPath)
			if err != nil {
				return err
			}
			branch, err := r.Head()
			if err != nil {
				return err
			}
			gitRepos = append(gitRepos, &repo.Repo{Name: filepath.Base(repoPath), Path: repoPath, Branch: branch.Name().Short()})
			return fs.SkipDir
		}

		return nil
	})
	return gitRepos, nil
}
