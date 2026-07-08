package router

import (
	page "git.shi.foo/pages/repos"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("")

	urls.Fallback(urls.Get, "/:owner/:repo", page.Show, "repo.show")
}
