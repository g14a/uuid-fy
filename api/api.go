package api

import (
	"encoding/json"
	"net/http"
	"time"
	"uuid-fy/controller"
	"uuid-fy/jwtauth"
	"uuid-fy/models"
)

func CreatePerson(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()

	var person models.PersonModel
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	result, err := controller.CreatePerson(person)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, result.(string))
	}

	respondWithJSON(w, http.StatusOK, result)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var person models.UpdatePersonModel

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	result, err := controller.UpdatePerson(person.Name, person)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, result.(string))
	}

	respondWithJSON(w, http.StatusOK, result)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var person models.PersonModel

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	result, err := controller.GetPerson(person.Name)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, result.(string))
	}

	respondWithJSON(w, http.StatusOK, result)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
	results, err := controller.GetAll()
	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	
	respondWithJSON(w, http.StatusOK, results)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds jwtauth.Credentials
	
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	expectedPassword, ok := Users[creds.Username]
	
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	tokenString, expirationTime, err := jwtauth.JwtToken(creds.Username)
	
	http.SetCookie(w, &http.Cookie{
		Name:       "token",
		Value:      tokenString,
		Expires:   expirationTime,
	})
}

func respondWithError(w http.ResponseWriter, httpCode int, message string) {
	respondWithJSON(w, httpCode, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, httpCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, _ = w.Write(response)
}

// Username, password map
var Users = map[string]string {
	"gowtham": "clear",
}
