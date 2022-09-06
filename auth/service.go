package auth

import (
	"bwastartup/helper"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func NewJWTService() *jwtService {
	return &jwtService{}
}

var secretKey = helper.GoDotEnvVariable("SECRET_KEY")

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwtService{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "error", err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid Token")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
