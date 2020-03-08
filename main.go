package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uuid-fy/api"
)

func main() {
	r := mux.NewRouter()
	
	r.Handle("/users/create", api.IsAuthorized(api.CreatePerson)).Methods("POST")
	r.Handle("/users/update", api.IsAuthorized(api.UpdatePerson)).Methods("POST")
	r.Handle("/users", api.IsAuthorized(api.GetPerson)).Methods("GET")
	r.Handle("/getall", api.IsAuthorized(api.GetAllUsers)).Methods("GET")
	r.HandleFunc("/signin", api.Signin).Methods("POST")
	r.HandleFunc("/signup", api.Signup).Methods("POST")
	log.Fatal(http.ListenAndServeTLS(":8000", "server.crt", "server.key", r))
	
}
