package requests

type LocalLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleLoginRequest struct {
	Email string `json:"email"`
}
