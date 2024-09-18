package dto

type RequestUpdateUser struct {
	Name string `json:"name,omitempty" validate:"required,min=3"`
}

type RequestUpdatePassword struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}