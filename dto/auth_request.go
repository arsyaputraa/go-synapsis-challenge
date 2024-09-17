package dto

type RequestRegister struct {
	Name     string `json:"name,omitempty" validate:"required,min=3"`
	Email    string `json:"email,omitempty"  validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type RequestLogin struct {
	Email    string `json:"email,omitempty"  validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}