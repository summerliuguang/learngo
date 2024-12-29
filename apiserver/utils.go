package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("")
var jwtIssuer = ""
var jwtExpire = time.Hour * 24
var turnstileKey = ""

type Claims struct {
	Username string
	jwt.StandardClaims
}

// SetJWTKey set the jwt key
func init() {
	jwtKey = []byte(os.Getenv("JWT_KEY"))
	jwtIssuer = os.Getenv("JWT_ISSUER")
	turnstileKey = os.Getenv("TURNSTILE_KEY")
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

func VaildTurnstile(turnsfileData TurnstileVerify) bool {
	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	turnsfileData.Secret = turnstileKey
	log.Println("token:", turnsfileData.Response)
	log.Println("secret:", turnsfileData.Secret)
	jsonData, err := json.Marshal(turnsfileData)
	if err != nil {
		log.Println("Error in Marshal", err)
		return false
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error in NewRequest", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error in Do", err)
		return false
	}
	defer resp.Body.Close()
	var result TurnstileResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error in Decode", err)
		return false
	}
	log.Println("Turnstile Response:", result.Success)
	log.Println("Turnstile Chllenge_ts:", result.Chllenge_ts)
	log.Println("Turnstile Hostname:", result.Hostname)
	return result.Success
}
