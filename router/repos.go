package router

import (
	controller "git.shi.foo/controllers/repos"
	page "git.shi.foo/pages/repos"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("repos")

	urls.Path(urls.Get, "/", page.Index, "index")
	urls.Path(urls.Get, "/import", page.ImportIndex, "import")
	urls.Path(urls.Post, "/import", controller.Import, "import.create")
}
