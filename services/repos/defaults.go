package repos

const (
	LogPrefix   = "Repos"
	ReposTitle  = "Repositories"
	ImportTitle = "Import repositories"

	SizeFormatMB = "%.1f MB"
	SizeFormatKB = "%.1f KB"
	SizeFormatB  = "%d B"
	DirSizeLabel = "-"

	PlainFormat = "<pre class=\"blob-code-plain\">%s</pre>"

	ImportMaxAttempts   = 3
	ImportEventsURLPath = "/import/events"
	DonePercent         = 100

	bytesPerKB       = 1 << 10
	bytesPerMB       = 1 << 20
	binarySniffLimit = 8000
)
