package pgfunc

import (
	"context"
	"fmt"
	"log"
	"sync"
	config2 "uuid-fy/config"
	"uuid-fy/models"

	"github.com/go-pg/pg/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	pgDB *pg.DB
	once sync.Once
)

func AddUserAuthData(user models.User) error {
	err := pgDB.Insert(&user)
	if err != nil {
		return err
	}

	return nil
}

// Check if User exists
func CheckUser(username, password string) bool {
	var m models.User
	ormQuery := pgDB.Model(&m)

	ormQuery = ormQuery.Where("username = ?", username)
	err := ormQuery.Select(&m)
	if err != nil {
		return false
	}

	if !CheckPasswordHash(password, m.Password) {
		return false
	}

	return true
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SetupPostgres() {
	once.Do(func() {
		db, err := PostgresClient()
		if err != nil {
			log.Println(err)
		}
		pgDB = db
	})
}

func PostgresClient() (*pg.DB, error) {
	config := config2.GetAppConfig()
	db := pg.Connect(&pg.Options{
		User:     config.PostgresConfig.Username,
		Password: config.PostgresConfig.Password,
		Database: config.PostgresConfig.DbName,
		Addr:     config.PostgresConfig.ServerURL,
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}
	
	db.AddQueryHook(dbLogger{})
	
	return db, nil
}

type dbLogger struct { }

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

func init() {
	SetupPostgres()
}
