package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uuid-fy/api"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/users/create", api.CreateUser).Methods("POST")
	r.HandleFunc("/users/update", api.UpdatePerson).Methods("POST")
	r.HandleFunc("/users/{name}", api.GetUser).Methods("GET")
	r.HandleFunc("/getall", api.GetAllUsers).Methods("GET")
	r.HandleFunc("/signin", api.Signin).Methods("POST")
	r.HandleFunc("/signup", api.Signup).Methods("POST")
	r.HandleFunc("/users/{username}/createcontactinfo", api.AddContactInfoToUser).Methods("POST")
	r.HandleFunc("/users/{username}/contactinfo",  api.GetContactInfoOfUser).Methods("GET")
	
	log.Fatal(http.ListenAndServeTLS(":8000", "server.crt", "server.key", r))
	
}
