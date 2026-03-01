package repo

type Repo struct {
	Name   string
	Path   string
	Branch string
}

func NewRepo(name string, path string, branch string) *Repo {
	return &Repo{Name: name, Path: path, Branch: branch}
}
