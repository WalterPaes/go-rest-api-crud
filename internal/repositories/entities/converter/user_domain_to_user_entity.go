package converter

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities"
)

func UserDomainToUserEntity(userDomain *domain.User) *entities.UserEntity {
	return &entities.UserEntity{
		Name:     userDomain.Name,
		Email:    userDomain.Email,
		Password: userDomain.Password,
	}
}
