package router

import (
	page "git.shi.foo/pages/repos"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("")

	urls.Fallback(urls.Get, "/:owner/:repo", page.Show, "repo.show")
	urls.Fallback(urls.Get, "/:owner/:repo/tree/*", page.Show, "repo.tree")
	urls.Fallback(urls.Get, "/:owner/:repo/blob/*", page.Blob, "repo.blob")
}
