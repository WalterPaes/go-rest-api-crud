package converter

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos"
)

func UserDomainToUserResponse(userDomain *domain.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:    userDomain.ID,
		Name:  userDomain.Name,
		Email: userDomain.Email,
	}
}
