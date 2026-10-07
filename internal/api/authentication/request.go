package authentication

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
