package repos

import "time"

const (
	LogPrefix   = "Repos"
	ReposTitle  = "Repositories"
	ImportTitle = "Import repositories"

	SizeFormatMB = "%.1f MB"
	SizeFormatKB = "%.1f KB"
	SizeFormatB  = "%d B"
	DirSizeLabel = "-"

	PlainFormat = "<pre class=\"blob-code-plain\">%s</pre>"

	ReadmePlainOpen  = "<pre class=\"readme-plain\">"
	ReadmePlainClose = "</pre>"

	CommitConcurrency = 8
	CommitsPageSize   = 40
	ImportMaxAttempts = 3

	InsightsCacheKey     = "repo:insights:%s/%s"
	InsightsTTL          = 15 * time.Minute
	LanguageLimit        = 6
	RailContributorLimit = 5
	CreatedDateFormat    = "2006-01-02"

	TabFiles            = "files"
	TabCommits          = "commits"
	TabBranches         = "branches"
	TabTags             = "tags"
	TabContributors     = "contributors"
	ImportEventsURLPath = "/import/events"
	DonePercent         = 100

	bytesPerKB       = 1 << 10
	bytesPerMB       = 1 << 20
	binarySniffLimit = 8000
)
