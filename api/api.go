package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uuid-fy/jwtauth"
	"uuid-fy/models"
	"uuid-fy/neofunc"
	"uuid-fy/neofunc/contact_info"
	"uuid-fy/pgfunc"
)

func CreateUser(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()

	var person models.UserModel
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	result, err := neofunc.CreateUser(person)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, result.(string))
	}

	respondWithJSON(w, http.StatusOK, result)
}

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

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	
	userName := params["name"]

	result, err := neofunc.GetUser(userName)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, result.(string))
	}

	respondWithJSON(w, http.StatusOK, result)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
	results, err := neofunc.GetAll()
	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	
	respondWithJSON(w, http.StatusOK, results)
}

func AddContactInfoToUser(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	username := params["name"]
	
	var contactNode models.ContactInfoModel
	
	err := json.NewDecoder(r.Body).Decode(&contactNode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	results, err := contact_info.CreateContactInfo(contactNode)
	results, err = contact_info.CreateRelationToContactNode(username, contactNode.Phone)
	
	if err != nil {
		log.Println(err)
	}
	
	respondWithJSON(w, http.StatusOK, results)
}

func GetContactInfoOfUser(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	username := params["username"]
	
	results, err := contact_info.GetContactInfoOfUser(username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	
	if pgfunc.CheckUser(creds.Username, creds.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	// get back jwt token to the client
	tokenString, expirationTime, err := jwtauth.JwtToken(creds.Username)
	
	http.SetCookie(w, &http.Cookie{
		Name:       "token",
		Value:      tokenString,
		Expires:   expirationTime,
	})
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
		c, err := r.Cookie("token")
		
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		tokenString := c.Value
		
		claims := &jwtauth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
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
		
		endpoint(w, r)
	})
}
