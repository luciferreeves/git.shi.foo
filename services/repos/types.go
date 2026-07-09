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

type EntryView struct {
	Type      string
	Name      string
	Size      int64
	SizeLabel string
	IsDir     bool
	URL       string
}

type Crumb struct {
	Name string
	URL  string
}

type CommitView struct {
	Short   string
	Message string
	Author  string
	When    time.Time
}

type ShowContext struct {
	Title         string
	Owner         string
	Name          string
	Private       bool
	Description   string
	DefaultBranch string
	Status        string
	CloneURL      string
	Ready         bool
	Path          string
	Crumbs        []Crumb
	LatestCommit  *CommitView
	Entries       []EntryView
	Importing     bool
	ImportPhase   string
	ImportPercent int
	EventsURL     string
	ReadmeName    string
	ReadmeHTML    string
}

type ImportStreamContext struct {
	Topic  string
	RepoID uint
}

type BlobContext struct {
	Title     string
	Owner     string
	Name      string
	Path      string
	Filename  string
	Crumbs    []Crumb
	Content   string
	IsBinary  bool
	SizeLabel string
}
