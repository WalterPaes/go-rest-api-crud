package converter

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities"
)

func UserEntityToUserDomain(userEntity entities.UserEntity) *domain.User {
	return &domain.User{
		ID:       userEntity.ID.Hex(),
		Name:     userEntity.Name,
		Email:    userEntity.Email,
		Password: userEntity.Password,
	}
}
