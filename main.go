package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"uuid-fy/api"
)

func main() {
	r := mux.NewRouter()
	
	go func() {
		time.Sleep()
	}()
	
	r.Handle("/users/create", api.IsAuthorized(api.CreatePerson)).Methods("POST")
	r.Handle("/users/update", api.IsAuthorized(api.UpdatePerson)).Methods("POST")
	r.Handle("/users", api.IsAuthorized(api.GetPerson)).Methods("GET")
	r.Handle("/getall", api.IsAuthorized(api.GetAllUsers)).Methods("GET")
	r.HandleFunc("/signin", api.Signin).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
	
}
