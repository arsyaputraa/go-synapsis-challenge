package utils

import (
	"time"

	"github.com/arsyaputraa/go-synapsis-challenge/pkg/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var jwtSecret = []byte(config.Config("JWT_SECRET_KEY"))

func GenerateJWT(userID uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func GetJWTSecret() []byte {
	return jwtSecret
}
