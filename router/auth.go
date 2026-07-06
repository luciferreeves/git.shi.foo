package router

import (
	controller "git.shi.foo/controllers/auth"
	page "git.shi.foo/pages/auth"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("auth")

	urls.Path(urls.Get, "/login", page.Login, "login")
	urls.Path(urls.Get, "/callback", controller.Callback, "callback")
	urls.Path(urls.Post, "/logout", controller.Logout, "logout")
}
