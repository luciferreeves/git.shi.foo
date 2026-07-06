package invites

const (
	TokenGenerationLog = "Failed to generate invite token: %v"
	InviteCreateLog    = "Failed to create invitation: %v"
	InviteEmailLog     = "Failed to send invite email: %v"
	InviteLookupLog    = "Failed to look up invitation: %v"
	InviteAcceptLog    = "Failed to accept invitation: %v"

	MissingFields      = "Email and username are required."
	InviteCreateFailed = "Could not create the invitation."
	InviteEmailFailed  = "Could not send the invitation email."
	MissingToken       = "Missing invitation token."
	InviteNotFound     = "Invitation not found."
	InviteAlreadyUsed  = "This invitation has already been used."
	InviteAcceptFailed = "Could not accept the invitation."

	InviteEmailSubject = "You're invited to git.shi.foo"
	InviteEmailBody    = "Hi %s,\n\nYou've been invited to git.shi.foo. Accept your invitation:\n\n%s\n\nThen sign in with GitHub.\n"
)
