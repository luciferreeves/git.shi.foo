package router

import (
	controller "git.shi.foo/controllers/sync"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("webhooks")

	urls.Path(urls.Post, "/github", controller.Webhook, "github")
}
