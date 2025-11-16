package register

import "github.com/dgrijalva/jwt-go"

// type RegisterRequest struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// LoginRequest Структура для входа
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims Структура для JWT токена
type Claims struct {
	Username           string `json:"username"`
	jwt.StandardClaims        // Данное поле нужно для правильной генерации JWT
}
