package validation

type ImagesMarkDown struct {
	Data Images `json:"data" validate:"required"`
}
type Images struct {
	ID    int   `json:"id" validate:"required,min=1"`
	Value Value `json:"value" validate:"required"`
}
type Value struct {
	FullName  string `json:"full_name" validate:"required,min=1"`
	Branch    string `json:"branch" validate:"required,min=1,max=15"`
	OwnerName string `json:"ownerName" validate:"required,min=1"`
}
