package router

import (
	controller "git.shi.foo/controllers/settings"
	page "git.shi.foo/pages/settings"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("settings")

	urls.Path(urls.Get, "/", page.Index, "index")
	urls.Path(urls.Post, "/tokens", controller.CreateToken, "tokens.create")
	urls.Path(urls.Post, "/tokens/:id/delete", controller.RevokeToken, "tokens.delete")
	urls.Path(urls.Post, "/keys", controller.AddKey, "keys.create")
	urls.Path(urls.Post, "/keys/import", controller.ImportKeys, "keys.import")
	urls.Path(urls.Post, "/keys/:id/delete", controller.RemoveKey, "keys.delete")
}
