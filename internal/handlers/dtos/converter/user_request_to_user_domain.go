package converter

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos"
)

func UserRequestToUserDomain(userRequest dtos.UserRequest) *domain.User {
	return &domain.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}
}
