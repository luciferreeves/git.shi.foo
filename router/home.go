package router

import (
	page "git.shi.foo/pages/home"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("")

	urls.Path(urls.Get, "/", page.Index, "home")
}
