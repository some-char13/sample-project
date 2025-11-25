package jwt

import (
	"errors"
	"sample_project/internal/conf"
	"sample_project/internal/model/register"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret               = []byte(conf.Load().JWTSecret)
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
)

type Worker struct{}

func New() *Worker {
	return &Worker{}
}

func (w *Worker) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &register.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "monitor-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*register.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &register.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*register.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
