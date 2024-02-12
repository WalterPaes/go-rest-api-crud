package converter

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos"
)

func LoginRequestToUserDomain(loginRequest dtos.LoginRequest) *domain.User {
	return &domain.User{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}
}
