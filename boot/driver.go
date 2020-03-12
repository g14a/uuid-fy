package boot

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"uuid-fy/config"
)

var Driver neo4j.Driver

func SetupDriver(config config.AppConfig) error {
	driver, err := neo4j.NewDriver(config.Neo4jConfig.ServerURL,
		neo4j.BasicAuth(config.Neo4jConfig.Username, config.Neo4jConfig.Password, ""),
		func(config *neo4j.Config) {
		config.Encrypted = false
	})

	if err != nil {
		return err
	}

	Driver = driver

	return nil
}

func GetWriteSession() (neo4j.Session, error) {

	session, err := Driver.Session(neo4j.AccessModeWrite)

	if err != nil {
		log.Println(err)
	}

	return session, err
}

func GetReadSession() (neo4j.Session, error) {

	session, err := Driver.Session(neo4j.AccessModeRead)

	if err != nil {
		log.Println(err)
	}

	return session, err
}

func init() {
	config := config.GetAppConfig()
	SetupDriver(*config)
}