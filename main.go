package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signup", SignUp)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/logout", LogOut)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
