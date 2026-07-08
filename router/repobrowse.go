package router

import (
	controller "git.shi.foo/controllers/repos"
	page "git.shi.foo/pages/repos"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("")

	urls.Fallback(urls.Get, "/:owner/:repo", page.Show, "repo.show")
	urls.Fallback(urls.Get, "/:owner/:repo/tree/*", page.Show, "repo.tree")
	urls.Fallback(urls.Get, "/:owner/:repo/blob/*", page.Blob, "repo.blob")
	urls.Fallback(urls.Get, "/:owner/:repo/import/events", page.ImportEvents, "repo.import.events")
	urls.Fallback(urls.Post, "/:owner/:repo/import/retry", controller.Retry, "repo.import.retry")
}
