package gitserve

const (
	LogPrefix = "GitServe"

	NoCache   = "no-cache, max-age=0, must-revalidate"
	AuthRealm = "Basic realm=\"git.shi.foo\""

	BasicPrefix       = "Basic "
	CredentialColon   = ":"
	GzipEncoding      = "gzip"
	AdvertisementType = "application/x-%s-advertisement"
	ResultType        = "application/x-%s-result"
)
