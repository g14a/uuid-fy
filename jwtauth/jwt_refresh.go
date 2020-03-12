package jwtauth

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

func RefreshJWT(tokenStr string) string {
	
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (i interface{}, err error) {
		return JwtKey, nil
	})
	
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return ""
		}
		return ""
	}
	
	if !token.Valid {
		return ""
	}
	
	// check if old token is about to expire only under 30 seconds
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30 * time.Second {
		return tokenStr
	}
	
	// Create new token with renewed expiration time
	expirationTime := time.Now().Add(1 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = newToken.SignedString(JwtKey)
	
	if err != nil {
		log.Println(err)
		return ""
	}
	
	return tokenStr
	
}