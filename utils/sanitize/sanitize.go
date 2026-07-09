package sanitize

import "github.com/microcosm-cc/bluemonday"

var policy = buildPolicy()

func buildPolicy() *bluemonday.Policy {
	sanitizer := bluemonday.UGCPolicy()
	sanitizer.AllowAttrs("class").Globally()
	sanitizer.AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6")
	return sanitizer
}

func HTML(raw []byte) string {
	return string(policy.SanitizeBytes(raw))
}
