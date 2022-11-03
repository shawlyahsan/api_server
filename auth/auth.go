package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var user string
var pass string

func BasicAuthentication(handle_func http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user = os.Getenv("api_user")
		pass = os.Getenv("api_pass")

		username, password, ok := r.BasicAuth()

		if ok == false {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if username == user && password == pass {
			handle_func.ServeHTTP(w, r)
		} else {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

	}
}

func Get_Token() (string, error) {
	var signingKey = []byte("samplesecretkey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     user,
		"exp_time": time.Now().Add(time.Minute * 15).Unix(),
	})

	token_string, err := token.SignedString(signingKey)
	return token_string, err
}

func verifyToken(token_string string) (jwt.Claims, error) {
	var signingKey = []byte("samplesecretkey")
	token, err := jwt.Parse(token_string, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	} else {
		return token.Claims, nil
	}
}

func Is_Authorized(handle_func http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token_string := r.Header.Get("Authorization")
		if len(token_string) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization Failed"))
			return
		}

		token_string = strings.Replace(token_string, "Bearer ", "", 1)
		claims, err := verifyToken(token_string)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}

		name := claims.(jwt.MapClaims)["name"].(string)
		w.Header().Set("name", name)

		handle_func.ServeHTTP(w, r)

	}
}
