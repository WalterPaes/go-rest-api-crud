package dtos

type UserRequest struct {
	Name     string `json:"name" binding:"required,min=4,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,containsany=!@#$%*"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersListResponse struct {
	Users        []UserResponse `json:"users"`
	CurrentPage  int            `json:"current_page"`
	TotalPerPage int            `json:"total_per_page"`
}
