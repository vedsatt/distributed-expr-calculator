package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	HmacSampleSecret string = "an7DkUH?L8iClxbVj5JZdbRVO2M$1Jc~D6CXsL@4"
)

func Generate(username string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": username,
		"nbf":  now.Unix(),                       // когда станет валидным
		"exp":  now.Add(10 * time.Minute).Unix(), // когда перестанет быть валидным
		"iat":  now.Unix(),                       // время создания
	})

	tokenString, err := token.SignedString([]byte(HmacSampleSecret))
	if err != nil {
		log.Printf("jwt generation error: %v", err)
		return "", err
	}

	return tokenString, nil
}

func Verify(tokenString string) bool {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(HmacSampleSecret), nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	if !tokenFromString.Valid {
		log.Println("invalid token")
		return false
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		log.Println("user name: ", claims["name"])
		return true
	}

	return false
}
