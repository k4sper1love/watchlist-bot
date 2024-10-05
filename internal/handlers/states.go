package handlers

const (
	RegistrationStateAwatingUsername = "registration_state_awaiting_username"
	RegistrationStateAwatingEmail    = "registration_state_awaiting_email"
	RegistrationStateAwatingPassword = "registration_state_awaiting_password"

	LoginStateAwaitingEmail    = "login_state_awaiting_email"
	LoginStateAwaitingPassword = "login_state_awaiting_password"

	LogoutStateAwaitingConfirm = "logout_state_awaiting_confirm"
)
