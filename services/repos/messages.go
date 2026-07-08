package repos

const (
	TokenLog      = "Failed to get access token: %v"
	ReposFetchLog = "Failed to fetch GitHub repos: %v"
	ListLog       = "Failed to list repos: %v"
	MetadataLog   = "Failed to fetch repo metadata: %v"
	CreateLog     = "Failed to create repo row: %v"
	MirrorLog     = "Failed to mirror repo: %v"
	HighlightLog  = "Failed to highlight file: %v"

	AuthRequired     = "Sign in required."
	ReposFetchFailed = "Could not fetch your GitHub repositories."
	ListFailed       = "Could not load repositories."
	ImportFailed     = "Could not import the repository."
	AlreadyImported  = "That repository is already imported."
	RepoNotFound     = "Repository not found."
	FileNotFound     = "File not found."
)
