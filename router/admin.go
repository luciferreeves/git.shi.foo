package router

import (
	controller "git.shi.foo/controllers/admin"
	page "git.shi.foo/pages/admin"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("admin")

	urls.Path(urls.Get, "/", page.Index, "index")
	urls.Path(urls.Post, "/invites", controller.CreateInvite, "invites")
}
