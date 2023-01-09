package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//"github.com/dgrijalva/jwt-go"
	jwt "github.com/dgrijalva/jwt-go"
	//"github.com/golang-jwt/jwt"
	//"github.com/golang-jwt/jwt"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "shabana"
	claims["aud"] = "billing.jwt.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(10 * time.Minute)
	//	claims["exp"] = time.Time.Add(time.Minute * 1).Unix()
	tokenString, err := token.SignedString(MySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	fmt.Println(validToken)
	if err != nil {
		fmt.Println("failed to generate token")
	}
	fmt.Fprintln(w, string(validToken))
}

func handleRequests() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
