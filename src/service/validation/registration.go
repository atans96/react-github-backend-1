package validation

type Registration struct {
	ClientId     string `json:"client_id" validate:"required,min=4,max=16,regexp=^[a-zA-Z0-9]*$"`
	RedirectURL  string `json:"redirect_uri" validate:"required,min=4,max=16,url"`
	ClientSecret string `json:"client_secret" validate:"required,min=1,max=16"`
	TokenCode    string `json:"code" validate:"required,min=1,max=16"`
}
