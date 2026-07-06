package admin

import "time"

type UserView struct {
	Login     string
	Email     string
	Admin     bool
	Enabled   bool
	CreatedAt time.Time
}

type InvitationView struct {
	ID        uint
	Email     string
	Username  string
	Status    string
	CreatedAt time.Time
}

type IndexContext struct {
	Title       string
	Users       []UserView
	Invitations []InvitationView
}
