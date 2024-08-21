package jwtoken

import (
	"errors"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	secretkey  = "very-secret-key"
	ExpiresKey = 72
)

// NewToken генерируем JWToken
func NewToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	if user.ID > 0 && user.Email != "" {
		claims["uid"] = user.ID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * ExpiresKey).Unix()
	} else {
		return "", errors.New("can't create JWT, invalid user id or login")
	}

	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

// проверка валидности токена
func GetJWTokenUserID(tokenString string) int {

	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretkey), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}

	userID := claims["uid"].(float64)

	return int(userID) //id
}
