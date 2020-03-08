package models

type UserModel struct {
	UUID string `json:"id"`
	Name string	`json:"name"`
}

type ContactInfoModel struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
	Address string `json:"address"`
}

// Auth user in Postgres
type User struct {
	Username string `pg:"username"`
	Password string `pg:"password"`
	Email    string `pg:"email"`
	tableName struct{}  `pg:"users"`
}

