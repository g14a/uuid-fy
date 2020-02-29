package boot

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"uuid-fy/config"
)

func SetupDriver(config config.AppConfig) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(config.Neo4jConfig.ServerURL,
		neo4j.BasicAuth(config.Neo4jConfig.Username, config.Neo4jConfig.Password, ""),
		func(config *neo4j.Config) {
		config.Encrypted = false
	})

	if err != nil {
		return nil, err
	}

	return driver, nil
}

func GetWriteSession() (neo4j.Session, error) {
	appConfig := config.GetAppConfig()
	driver, err := SetupDriver(*appConfig)
	if err != nil {
		log.Println(err)
	}

	session, err := driver.Session(neo4j.AccessModeWrite)

	if err != nil {
		log.Println(err)
	}

	return session, err
}

func GetReadSession() (neo4j.Session, error) {
	appConfig := config.GetAppConfig()
	driver, err := SetupDriver(*appConfig)
	if err != nil {
		log.Println(err)
	}

	session, err := driver.Session(neo4j.AccessModeRead)

	if err != nil {
		log.Println(err)
	}

	return session, err
}