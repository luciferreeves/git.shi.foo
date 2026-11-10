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
	Type          string
	Name          string
	Size          int64
	SizeLabel     string
	IsDir         bool
	URL           string
	CommitMessage string
	CommitAuthor  string
	CommitWhen    time.Time
}

type Crumb struct {
	Name string
	URL  string
}

type RefOption struct {
	Name   string
	URL    string
	Active bool
	IsTag  bool
}

type CommitView struct {
	Short   string
	Message string
	Author  string
	When    time.Time
}

type RefView struct {
	Name    string
	Short   string
	Message string
	Author  string
	URL     string
	When    time.Time
}

type ContributorView struct {
	Name  string
	Email string
	Count int
}

type LanguageView struct {
	Name    string
	Percent string
	Width   string
	Color   string
}

type BranchesContext struct {
	Title         string
	Owner         string
	Name          string
	Private       bool
	Description   string
	DefaultBranch string
	Status        string
	CloneURL      string
	SSHCloneURL   string
	ActiveTab     string
	Branches      []RefView
}

type TagsContext struct {
	Title         string
	Owner         string
	Name          string
	Private       bool
	Description   string
	DefaultBranch string
	Status        string
	CloneURL      string
	SSHCloneURL   string
	ActiveTab     string
	Tags          []RefView
}

type ContributorsContext struct {
	Title         string
	Owner         string
	Name          string
	Private       bool
	Description   string
	DefaultBranch string
	Status        string
	CloneURL      string
	SSHCloneURL   string
	ActiveTab     string
	Contributors  []ContributorView
}

type CommitsContext struct {
	Title         string
	Owner         string
	Name          string
	Private       bool
	Description   string
	DefaultBranch string
	Status        string
	CloneURL      string
	SSHCloneURL   string
	ActiveTab     string
	Ref           string
	RefBranches   []RefOption
	RefTags       []RefOption
	Commits       []CommitView
	Page          int
	HasPrev       bool
	HasNext       bool
	PrevURL       string
	NextURL       string
}

type ShowContext struct {
	Title            string
	Owner            string
	Name             string
	Private          bool
	Description      string
	DefaultBranch    string
	Status           string
	CloneURL         string
	SSHCloneURL      string
	ActiveTab        string
	Ready            bool
	Ref              string
	RefBranches      []RefOption
	RefTags          []RefOption
	RootURL          string
	Path             string
	Crumbs           []Crumb
	LatestCommit     *CommitView
	Entries          []EntryView
	Importing        bool
	ImportPhase      string
	ImportPercent    int
	EventsURL        string
	ReadmeName       string
	ReadmeHTML       string
	Topics           []string
	Homepage         string
	CreatedLabel     string
	Languages        []LanguageView
	RailContributors []ContributorView
}

type ImportStreamContext struct {
	Topic  string
	RepoID uint
}

type BlobContext struct {
	Title         string
	Owner         string
	Name          string
	Private       bool
	Description   string
	DefaultBranch string
	Status        string
	CloneURL      string
	SSHCloneURL   string
	ActiveTab     string
	Ref           string
	RefBranches   []RefOption
	RefTags       []RefOption
	RootURL       string
	Path          string
	Filename      string
	Crumbs        []Crumb
	Content       string
	IsBinary      bool
	SizeLabel     string
}
