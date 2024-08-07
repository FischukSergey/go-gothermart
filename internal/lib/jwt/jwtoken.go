package jwtoken

import (
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	secretkey = "very-secret-key"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
//type Claims struct {
//	jwt.RegisteredClaims
//	Uid string `json:"uid"`
//	Sub string `json:"sub"`
//}

// NewToken генерируем JWToken
func NewToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

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
	//id, err := strconv.Atoi(userID)
	//if err != nil {
	//	return -1
	//}

	return int(userID) //id
}
