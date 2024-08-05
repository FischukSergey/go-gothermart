package jwt

import (
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"uid": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte("very-secret-key"))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
