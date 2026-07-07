package repos

import "time"

type GitHubRepoView struct {
	Owner       string
	Name        string
	FullName    string
	Private     bool
	Description string
	Mirrored    bool
}

type ImportContext struct {
	Title string
	Repos []GitHubRepoView
}

type RepoView struct {
	Owner       string
	Name        string
	Private     bool
	Description string
	Status      string
	UpdatedAt   time.Time
}

type IndexContext struct {
	Title string
	Repos []RepoView
}
