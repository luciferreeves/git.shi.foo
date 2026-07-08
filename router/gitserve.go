package router

import (
	controller "git.shi.foo/controllers/gitserve"
	"git.shi.foo/utils/urls"
)

func init() {
	urls.SetNamespace("")

	urls.Fallback(urls.Get, "/:owner/:repo/info/refs", controller.InfoRefs, "git.inforefs")
	urls.Fallback(urls.Post, "/:owner/:repo/git-upload-pack", controller.UploadPack, "git.uploadpack")
	urls.Fallback(urls.Post, "/:owner/:repo/git-receive-pack", controller.ReceivePack, "git.receivepack")
}
