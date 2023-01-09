package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
)

//adds the secret key to the my jwt token
var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

//handler function has 2things(req,resp),req is smg the user sends to it
// and func receives it(hence its a ptr)and resp is smg sent back from 
//this func to the user,resp is the msg displayed
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret information")
}

//to check the calid token,the function returns the handler
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

					return nil, fmt.Errorf("invalid signing method ")
				}
				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

				if !checkAudience {
					return nil, fmt.Errorf("invalid auf")
				}
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)

				//checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf("invalid iss")
				}

				return MySigningKey, nil

			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "no authorization token provided")
		}
	})
}
//this fucnc basically checks that our server is authorised and runs on port 9001
//the requests function routes "route" to home page wherein the func checks the token is 
//authorised or not
func handleRequests() {
	http.Handle("/", isAuthorized(homePage))
	log.Fatal((http.ListenAndServe(":9001", nil)))
}

//Control of the program comes first here
func main() {
	fmt.Printf("server")
	handleRequests()
}
