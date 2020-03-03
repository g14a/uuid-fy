package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uuid-fy/api"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/users/create", api.CreatePerson).Methods("POST")
	r.HandleFunc("/users/update", api.UpdatePerson).Methods("POST")
	r.HandleFunc("/users", api.GetPerson).Methods("GET")
	r.HandleFunc("/getall", api.GetAllUsers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
