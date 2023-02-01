package helper

import (
	"campyuk-api/config"

	"github.com/golang-jwt/jwt"
)

func ExtractToken(t interface{}) (uint, string) {
	user := t.(*jwt.Token)
	var (
		userID uint
		role   string
	)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		switch claims["userID"].(type) {
		case float64:
			userID = uint(claims["userID"].(float64))
		case int:
			userID = claims["userID"].(uint)
		}
		role = claims["role"].(string)
	}
	return userID, role
}

func GenerateJWT(id int, role string) (string, interface{}) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = id
	claims["role"] = role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	useToken, _ := token.SignedString([]byte(config.JWT_KEY))
	return useToken, token
}
