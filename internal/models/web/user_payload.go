package web

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	GoogleID string `json:"google_id" `
}
