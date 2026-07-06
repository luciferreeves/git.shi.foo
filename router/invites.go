package router

import (
	controller "git.shi.foo/controllers/invites"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("invites")

	urls.Path(urls.Get, "/accept", controller.Accept, "accept")
}
