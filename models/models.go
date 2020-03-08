package models

type PersonModel struct {
	UUID string `json:"id"`
	Name string	`json:"name"`
	Email string `json:"email"`
	DateOfBirth string `json:"dob"`
}

type UpdatePersonModel struct {
	Name string	`json:"name"`
	Email string `json:"email"`
	DateOfBirth string `json:"dob"`
}

// Auth user in Postgres
type User struct {
	Username string `pg:"username"`
	Password string `pg:"password"`
	Email    string `pg:"email"`
	tableName struct{}  `pg:"users"`
}

