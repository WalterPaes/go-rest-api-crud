package services

import "github.com/WalterPaes/go-rest-api-crud/internal/domain"

type UserService interface{}

type userSvc struct{}

func NewUserService() *userSvc {
	return &userSvc{}
}

func (us *userSvc) CreateUser(user domain.User) (*domain.User, error) {
	return nil, nil
}
