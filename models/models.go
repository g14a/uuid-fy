package models

type UserModel struct {
	UUID string `json:"id"`
	Name string	`json:"username"`
}

type ContactInfoModel struct {
	Name string `json:"name"`
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

type EducationInfoModel struct {
	RootID string `json:"rootuid"`
	PrimarySchool string `json:"primary"`
	SecondarySchool string `json:"secondary"`
	University string `json:"university"`
}