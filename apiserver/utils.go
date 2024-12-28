package apiserver

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("")
var jwtIssuer = ""
var jwtExpire = time.Hour * 24

type Claims struct {
	Username string
	jwt.StandardClaims
}

// SetJWTKey set the jwt key
func init() {
	jwtKey = []byte(os.Getenv("JWT_KEY"))
	jwtIssuer = os.Getenv("JWT_ISSUER")
	if len(jwtKey) == 0 || len(jwtIssuer) == 0 {
		log.Fatal("JWT_KEY and JWT_ISSUER must be set")
	}
}

func GenerateJWT(username string) (string, error) {
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpire).Unix(),
			Issuer:    jwtIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidJWT(signedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	log.Println("1 ", token, err)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
