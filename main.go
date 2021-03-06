package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uuid-fy/api"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/users/getall", api.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/update", api.UpdatePerson).Methods("POST")
	r.HandleFunc("/users/{username}", api.GetUserUUID).Methods("GET")
	r.HandleFunc("/users/signin", api.Signin).Methods("POST")
	r.HandleFunc("/users/signup", api.Signup).Methods("POST")
	
	// Contact Info
	r.HandleFunc("/users/{username}/addcontactinfo", api.AddContactInfoToUser).Methods("POST")
	r.HandleFunc("/users/{username}/contactinfo",  api.GetContactInfoOfUser).Methods("GET")
	r.HandleFunc("/users/{username}/deletecontactinfo", api.DeleteContactInfoOfUser).Methods("DELETE")
	
	// Education Info
	r.HandleFunc("/users/{username}/addeducationinfo", api.AddEducationInfoToUser).Methods("POST")
	r.HandleFunc("/users/{username}/educationinfo", api.GetEducationInfoOfUser).Methods("GET")
	
	// Health Info
	r.HandleFunc("/users/{username}/addhealthinfo", api.AddHealthInfoToUser).Methods("POST")
	r.HandleFunc("/users/{username}/healthinfo", api.GetHealthInfoOfUser).Methods("GET")
	
	log.Fatal(http.ListenAndServe(":8000", r))
	
}
