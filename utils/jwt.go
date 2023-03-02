package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type MyCustomClaims struct {
	UserId    uint
	SessionId uint
	jwt.RegisteredClaims
}

func GenerateAuthToken(UserId, SessionId uint) string {
	claims := MyCustomClaims{
		UserId,
		SessionId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte("Hello blockchain")
	tokenString, err := token.SignedString(key)
	if err != nil {
		logrus.Panicf("not able to generate the token %v", err)
	}
	return tokenString
}

type SessionBody struct {
	UserId    uint
	SessionId uint
}

var SessionData = SessionBody{}

func ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		hmacSampleSecret := []byte("Hello blockchain")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		SessionData.UserId = uint(claims["UserId"].(float64))
		SessionData.SessionId = uint(claims["SessionId"].(float64))
		return nil
	} else {
		return err
	}

}
