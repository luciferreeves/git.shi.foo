package admin

type CreateInviteRequest struct {
	Email    string `form:"email"`
	Username string `form:"username"`
}
