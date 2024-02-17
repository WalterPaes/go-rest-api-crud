package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/mocks"
	"github.com/WalterPaes/go-rest-api-crud/pkg/jwt"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
)

var (
	claims = map[string]any{
		"id":    responseUser.ID,
		"email": responseUser.Email,
		"name":  responseUser.Name,
	}

	token = "token_test"

	unauthorizedError = resterrors.NewUnauthorizedError(errInvalidCredentials)
)

func Test_loginSvc_LoginUser(t *testing.T) {
	type fields struct {
		userRepository repositories.UserRepository
		jwtAuth        jwt.JwtAuth
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr *resterrors.RestErr
	}{
		{
			name: "Should generate a jwt token when user login is successful",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					responseUser.EncryptPassword()

					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(responseUser, nil)
					return m
				}(),
				jwtAuth: func() jwt.JwtAuth {
					m := mocks.NewJwtAuth(t)
					m.On("GenerateToken", claims).
						Return(token, nil)
					return m
				}(),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    token,
			wantErr: nil,
		},
		{
			name: "Should return an unauthorized error when user login credentials are incorrect",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(nil, resterrors.NewNotFoundError("error"))
					return m
				}(),
				jwtAuth: mocks.NewJwtAuth(t),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    "",
			wantErr: unauthorizedError,
		},
		{
			name: "Should return an error when try find user by email",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(nil, internalServerError)
					return m
				}(),
				jwtAuth: mocks.NewJwtAuth(t),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    "",
			wantErr: internalServerError,
		},
		{
			name: "Should return an error when password is incorrect",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(&domain.User{}, nil)
					return m
				}(),
				jwtAuth: mocks.NewJwtAuth(t),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    "",
			wantErr: unauthorizedError,
		},
		{
			name: "Should return an error when password is incorrect",
			fields: fields{
				userRepository: func() repositories.UserRepository {
					m := mocks.NewUserRepository(t)
					m.On("FindUserByEmail", ctx, inputUser.Email).
						Return(responseUser, nil)
					return m
				}(),
				jwtAuth: func() jwt.JwtAuth {
					m := mocks.NewJwtAuth(t)
					m.On("GenerateToken", claims).
						Return("", internalServerError)
					return m
				}(),
			},
			args: args{
				ctx:  ctx,
				user: inputUser,
			},
			want:    "",
			wantErr: internalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLoginService(tt.fields.userRepository, tt.fields.jwtAuth)
			got, err := s.LoginUser(tt.args.ctx, tt.args.user)
			if got != tt.want {
				t.Errorf("loginSvc.LoginUser() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("loginSvc.LoginUser() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
