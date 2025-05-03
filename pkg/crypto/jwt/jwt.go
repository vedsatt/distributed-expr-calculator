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

func Generate(id int) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"nbf": now.Unix(),                       // когда станет валидным
		"exp": now.Add(10 * time.Minute).Unix(), // когда перестанет быть валидным
		"iat": now.Unix(),                       // время создания
	})

	tokenString, err := token.SignedString([]byte(HmacSampleSecret))
	if err != nil {
		log.Printf("jwt generation error: %v", err)
		return "", err
	}

	return tokenString, nil
}

func Verify(tokenString string) (bool, int) {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(HmacSampleSecret), nil
	})

	if err != nil {
		log.Println(err)
		return false, 0
	}

	if !tokenFromString.Valid {
		log.Println("invalid token")
		return false, 0
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		userID, ok := claims["id"].(float64)
		if !ok {
			log.Println("invalid user ID type")
			return false, 0
		}
		log.Println("user id: ", int(userID))
		return true, int(userID)
	}

	return false, 0
}
