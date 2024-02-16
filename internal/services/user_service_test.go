package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/mocks"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	logger.Init("debug", "stdout")
}

var (
	ctx = context.Background()

	userID = primitive.NewObjectID().Hex()

	inputUser = &domain.User{
		Name:     "First User",
		Email:    "firstuser@email.com",
		Password: "123456",
	}

	responseUser = &domain.User{
		ID:       userID,
		Name:     "First User",
		Email:    "firstuser@email.com",
		Password: "123456",
	}

	internalServerError = resterrors.NewInternalServerError("error")
)

func Test_userSvc_CreateUser(t *testing.T) {
	type fields struct {
		userRepository repositories.UserRepository
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr *resterrors.RestErr
	}{
		{
			name: "Should create an user without errors",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(nil, nil)

					m.On("CreateUser", ctx, inputUser).
						Return(responseUser, nil)
					return m
				}(),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    responseUser,
			wantErr: nil,
		},
		{
			name: "Should return an error when try call repository to create an user",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(nil, nil)

					m.On("CreateUser", ctx, inputUser).
						Return(nil, internalServerError)
					return m
				}(),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    nil,
			wantErr: internalServerError,
		},
		{
			name: "Should return an error when try call repository to find an user by email",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(nil, internalServerError)
					return m
				}(),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    nil,
			wantErr: internalServerError,
		},
		{
			name: "Should return an error when user is already registered",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(responseUser, resterrors.NewNotFoundError(errEmailAlreadyRegistered))
					return m
				}(),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    nil,
			wantErr: resterrors.NewBadRequestError(errEmailAlreadyRegistered),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(tt.fields.userRepository)

			got, err := s.CreateUser(tt.args.ctx, tt.args.user)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userSvc.CreateUser() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("userSvc.CreateUser() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userSvc_FindUserById(t *testing.T) {
	type fields struct {
		userRepository repositories.UserRepository
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr *resterrors.RestErr
	}{
		{
			name: "Should find an user without errors",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserById", ctx, userID).
						Return(responseUser, nil)
					return m
				}(),
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want:    responseUser,
			wantErr: nil,
		},
		{
			name: "Should return an error when try find user by email",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserById", ctx, userID).
						Return(nil, internalServerError)
					return m
				}(),
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want:    nil,
			wantErr: internalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(tt.fields.userRepository)

			got, err := s.FindUserById(tt.args.ctx, tt.args.userID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userSvc.FindUserById() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("userSvc.FindUserById() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
