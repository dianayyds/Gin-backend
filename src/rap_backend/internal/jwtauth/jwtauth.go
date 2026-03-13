package jwtauth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// JWT ...
type JWT struct {
	SigningKey []byte
}

// NewJWT ...
func NewJWT() *JWT {
	return &JWT{
		[]byte("rap_user_signkey"),
	}
}

// CustomClaims ...
type CustomClaims struct {
	UserId   int64  `json:"userid"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

// CreateToken ...
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParserToken ...
func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token不可用")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token过期")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("无效的token")
			} else {
				return nil, fmt.Errorf("token不可用")
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token无效")
}
