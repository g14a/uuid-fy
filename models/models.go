package models

import (
	"time"
)

type UserModel struct {
	UUID string `json:"id"`
	Username string	`json:"username"`
}

type ContactInfoModel struct {
	RootId string `json:"rootuid"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
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

type HealthInfoModel struct {
	RootID string `json:"rootuid"`
	BirthHospital string `json:"birthhospital"`
	BloodGroup string `json:"bloodgroup"`
}

type BlockChainModel struct {
	Id           int       `pg:"id,pk"`
	TimeStamp    time.Time `pg:"timestamp"`
	PreviousHash string    `pg:"previoushash"`
	Hash         string    `pg:"hash"`
	Data         NeoEvent  `pg:"data,type:json"`
	tableName    struct{}  `pg:"blockchain"`
}

type NeoEvent struct {
	EventType   string      `pg:"event"`
	DataPayload interface{} `pg:"data"`
	Message     string      `pg:"message"`
}
