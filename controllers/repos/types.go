package repos

type ImportRequest struct {
	Owner string `form:"owner"`
	Name  string `form:"name"`
}
