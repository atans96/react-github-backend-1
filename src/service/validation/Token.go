package validation

type Token struct {
	Token string `json:"token" validate:"required,min=10"`
}
