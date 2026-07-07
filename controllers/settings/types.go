package settings

type CreateTokenRequest struct {
	Label string `form:"label"`
}

type AddKeyRequest struct {
	Title string `form:"title"`
	Key   string `form:"key"`
}
