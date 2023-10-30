package userhttp

type GetUserByIDRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,ascii,min=3,max=32"`
}
