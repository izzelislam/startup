package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct{}

func NewService() *jwtService {
	return &jwtService{}
}

var SECRET_KEY = []byte("BWASTARTUP_s3cr3t_k3y")

func (s *jwtService) GenerateToken(userID int) (string, error) {
	paylod := jwt.MapClaims{}
	paylod["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, paylod)
	signedToken, err := token.SignedString(SECRET_KEY)

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}
