package appjwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/night1010/everhealth/config"
	"github.com/night1010/everhealth/entity"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Id     uint          `json:"id"`
	RoleId entity.RoleId `json:"role_id"`
}

type Jwt interface {
	GenerateToken(user *entity.User) (string, error)
	ValidateToken(tokenString string) (*entity.User, error)
}

type jwtImpl struct {
	secretKey []byte
}

func NewJwt() Jwt {
	jwtConfig := config.NewJwtConfig()
	return &jwtImpl{
		secretKey: []byte(jwtConfig.Secret),
	}
}

func (j *jwtImpl) GenerateToken(user *entity.User) (string, error) {
	jwtConfig := config.NewJwtConfig()
	appConfig := config.NewAppConfig()

	userId := strconv.Itoa(int(user.Id))
	duration := jwtConfig.ExpiryDuration

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    appConfig.Name,
			Subject:   userId,
		},
		Id:     user.Id,
		RoleId: user.RoleId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(j.secretKey)
	return signedString, err
}

func (j *jwtImpl) ValidateToken(tokenString string) (*entity.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok {
		user := &entity.User{
			Id:     claims.Id,
			RoleId: claims.RoleId,
		}
		return user, nil
	}
	return nil, errors.New("invalid claims type")
}
