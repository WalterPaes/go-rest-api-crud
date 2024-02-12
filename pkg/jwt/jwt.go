package jwt

import (
	"fmt"
	"time"

	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"github.com/golang-jwt/jwt"
)

type JwtAuth interface {
	GenerateToken(claims map[string]any) (string, *resterrors.RestErr)
}

type jwtAuth struct {
	secret  string
	expTime int
}

func NewJwtAuth(secret string, expTime int) *jwtAuth {
	return &jwtAuth{
		secret:  secret,
		expTime: expTime,
	}
}

func (a *jwtAuth) GenerateToken(claims map[string]any) (string, *resterrors.RestErr) {
	jwtClaims := jwt.MapClaims{"exp": time.Now().Add(time.Hour * time.Duration(a.expTime)).Unix()}

	for k, v := range claims {
		jwtClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", resterrors.NewInternalServerError(
			fmt.Sprintf("Error trying to generate jwt token: %s", err.Error()))
	}

	return tokenString, nil
}
