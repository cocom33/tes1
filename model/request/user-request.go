package request

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserUpdateRequest struct {
	Name    string `json:"name" validate:"required,min=3"`
	Address string `json:"address" validate:"required,min=5"`
	Phone   string `json:"phone"`
}

type UserEmailUpdateRequest struct {
	Email string `json:"email" validate:"required,email"`
}