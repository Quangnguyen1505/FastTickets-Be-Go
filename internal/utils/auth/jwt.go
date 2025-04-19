package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ntquang/ecommerce/global"
)

type PayLoadClams struct {
	jwt.RegisteredClaims
}

func generateTokenJWT(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.JWT.API_SERCERT_KEY))
}

func CreateToken(uuidToken string) (string, error) {
	//1. set time expiration
	timeEx := global.Config.JWT.JWT_EXPIRATION
	if timeEx == "" {
		timeEx = "1h"
	}

	expiration, err := time.ParseDuration(timeEx)
	if err != nil {
		return "", err
	}
	now := time.Now()
	expiresAt := jwt.NewNumericDate(now.Add(expiration))
	return generateTokenJWT(&PayLoadClams{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuidToken,
			Issuer:    "shopdevgo",
			ID:        uuid.New().String(),
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
}

func ParseJwtTokenSubject(token string, publickey string) (*jwt.RegisteredClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(JwtToken *jwt.Token) (interface{}, error) {
		return []byte(publickey), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.RegisteredClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
	// Parse the token with the correct claims type
}

// validate jwt token by subject
func VerifyToken(token string, publickey string) (*jwt.RegisteredClaims, error) {
	claims, err := ParseJwtTokenSubject(token, publickey)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
