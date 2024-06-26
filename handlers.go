package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var ListOfTokens = map[string]bool{}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func SignUp(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := users[credentials.Username]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("The username %s is already taken", credentials.Username)))
		return
	}

	if len(credentials.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("The password must have atleast 8 charcters/numbers")))
		return
	}

	users[credentials.Username] = credentials.Password
}

func Login(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[credentials.Username]
	if !ok || credentials.Password != expectedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, ok := ListOfTokens[tokenString]; !ok {
		ListOfTokens[tokenString] = true
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

}

func Home(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	if _, ok := ListOfTokens[tokenString]; !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
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
		delete(ListOfTokens, tokenString)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello %s", claims.Username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	if _, ok := ListOfTokens[tokenString]; !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
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
		delete(ListOfTokens, tokenString)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value

	delete(ListOfTokens, tokenString)
}
