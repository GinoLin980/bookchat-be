package dto

type UsersRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UsersResponse struct {
	Status   string `json:"status"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
