package jwtauth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func RefreshJWT(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	tokenString := c.Value
	
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return JwtKey, nil
	})
	
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	// check if old token is about to expire only under 30 seconds
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30 * time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// Create new token with renewed expiration time
	expirationTime := time.Now().Add(1 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = newToken.SignedString(JwtKey)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	http.SetCookie(w, &http.Cookie{
		Name:       "token",
		Value:      tokenString,
		Expires:    time.Time{},
	})
}