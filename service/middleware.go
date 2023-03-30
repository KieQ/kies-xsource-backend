package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

func ValidateToken(tokenStr string) (map[string]any, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid JWT method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Valid {
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			return mapClaims, nil
		}
		return nil, errors.New("claims is not map")
	} else {
		return nil, errors.New("validation has encountered unreachable code")
	}
}
