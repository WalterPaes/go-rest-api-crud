package converter

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
)

func UserRequestToUserDomain(userRequest dtos.UserRequest) (*domain.User, error) {
	user := &domain.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}
	err := user.EncryptPassword()
	if err != nil {
		logger.Error("Error when trying Encrypt Password", err)
		return nil, err
	}
	return user, nil
}
