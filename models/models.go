package models

type PersonModel struct {
	UUID string `json:"id"`
	Name string	`json:"name"`
	Email string `json:"email"`
	DateOfBirth string `json:"dob"`
}
