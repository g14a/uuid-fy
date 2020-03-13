package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"uuid-fy/jwtauth"
	"uuid-fy/models"
	"uuid-fy/neofunc"
	"uuid-fy/neofunc/contact_info"
	"uuid-fy/neofunc/education_info"
	"uuid-fy/pgfunc"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var contact models.ContactInfoModel

	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	result, err := contact_info.UpdateContactInfo(contact.Phone, contact)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, result.(string))
	}

	respondWithJSON(w, http.StatusOK, result)
}

func GetUserUUID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userName := params["username"]

	result, err := neofunc.GetUserUUID(userName)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, result.(string))
	}

	respondWithJSON(w, http.StatusOK, result)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
	results, err := neofunc.GetAll()

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	
	respondWithJSON(w, http.StatusOK, results)
}

func AddContactInfoToUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	var contactNode models.ContactInfoModel

	err := json.NewDecoder(r.Body).Decode(&contactNode)
	if err != nil {
		log.Println(err)
		return
	}

	results, err := contact_info.CreateContactInfo(contactNode)
	results, err = contact_info.CreateRelationToContactNode(username, contactNode.Phone)

	if err != nil {
		log.Println(err)
	}

	respondWithJSON(w, http.StatusOK, results)
}

func AddEducationInfoToUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	var educationNode models.EducationInfoModel

	err := json.NewDecoder(r.Body).Decode(&educationNode)
	if err != nil {
		log.Println(err)
		return
	}

	userUUID, err := neofunc.GetUserUUID(username)
	if err != nil {
		log.Println(err)
	}
	
	educationNode.RootID = userUUID.(string)

	results, err := education_info.CreateEducationInfo(educationNode)
	results, err = education_info.CreateRelationToEducationNode(educationNode.RootID)

	if err != nil {
		log.Println(err)
	}

	respondWithJSON(w, http.StatusOK, results)
}

func GetContactInfoOfUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	results, err := contact_info.GetContactInfoOfUser(username)

	if err != nil {
		log.Println(err)
	}

	respondWithJSON(w, http.StatusOK, results)
}

func GetEducationInfoOfUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	results, err := education_info.GetEducationInfoOfUser(username)

	if err != nil {
		log.Println(err)
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

	if !pgfunc.CheckUser(creds.Username, creds.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get back jwt token to the client
	tokenString, err := jwtauth.JwtToken(creds.Username)
	fmt.Println(tokenString)

	w.Header().Set("token", tokenString)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "Successful Sign in"})
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Password != "" && user.Username != "" {

		hashedPassword, err := pgfunc.HashPassword(user.Password)
		if err != nil {
			log.Println(err)
			return
		}

		user.Password = hashedPassword
		
		var userNode models.UserModel
		userNode.Username = user.Username

		_, err = neofunc.CreateUser(userNode)		

		err = pgfunc.AddUserAuthData(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
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

func IsAuthorized(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		// split by "bearer " and not just "bearer". Mind the extra space
		splitToken := strings.Split(tokenStr, "Bearer ")
		tokenStr = splitToken[1]

		if tokenStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		refreshToken := jwtauth.RefreshJWT(tokenStr)

		claims := &jwtauth.Claims{}
		token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (i interface{}, err error) {
			return jwtauth.JwtKey, nil
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("token", tokenStr)
		endpoint(w, r)
	})
}
