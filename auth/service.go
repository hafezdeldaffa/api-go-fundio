package auth

import (
	"bwastartup/helper"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct {
}

func NewJWTService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	var secretKey = helper.GoDotEnvVariable("SECRET_KEY")

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "error", err
	}

	return signedToken, nil
}
