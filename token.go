package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("your-secret-key") // Replace with your own secret key

// GenerateToken generates a JWT token with the provided claims.
func GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the claims if the token is valid.
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func userTokenGen(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (24 hours)
	}

	token, err := GenerateToken(claims)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return "", err
	}
	fmt.Println("Generated Token:", token)

	return token, nil
}

func validateUserToken(username string, r *http.Request) (bool, error) {
	token, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		fmt.Println("No cookie found")
		return false, err
	} else if err != nil {
		return false, err
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (24 hours)
	}

	tokenValue := token.Value

	claims, err = ValidateToken(tokenValue)
	if err != nil {
		fmt.Println("Error validating token:", err)
		return false, err
	}

	fmt.Println("Valid Token. Claims:", claims)
	return true, nil
}
