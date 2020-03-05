package jwtauth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"uuid-fy/config"
)

var JwtKey = []byte(config.GetAppConfig().AuthConfig.JwtToken)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func JwtToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(30 * time.Second)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		Username: username,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	
	if err != nil {
		return "", time.Time{}, err
	}
	
	return tokenString, expirationTime, nil
}