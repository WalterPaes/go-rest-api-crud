package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const errInvalidToken = "Invalid Token"

type JwtAuth interface {
	GenerateToken(claims map[string]any) (string, *resterrors.RestErr)
	VerifyTokenMiddleware(c *gin.Context)
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

func (a *jwtAuth) VerifyTokenMiddleware(c *gin.Context) {
	tokenValue := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")

	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(a.secret), nil
		}
		return nil, resterrors.NewBadRequestError(errInvalidToken)
	})
	if err != nil {
		errRest := resterrors.NewUnauthorizedError(errInvalidToken)
		c.JSON(errRest.HttpStatusCode, errRest)
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		errRest := resterrors.NewUnauthorizedError(errInvalidToken)
		c.JSON(errRest.HttpStatusCode, errRest)
		c.Abort()
		return
	}

	logger.Info(fmt.Sprintf("User authenticated: %+v", domain.User{
		ID:    claims["id"].(string),
		Name:  claims["name"].(string),
		Email: claims["email"].(string),
	}))
}
