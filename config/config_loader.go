package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

var appConfigInstance *AppConfig
var once sync.Once

// GetAppConfig returns an AppConfig instance
func GetAppConfig() *AppConfig {
	once.Do(func() {
		loadConfig()
	})

	return appConfigInstance
}

// loadConfig loads the data in the yaml file into a struct
// returns the app instance once if it is ready
func loadConfig() {
	yamlFile, err := ioutil.ReadFile("config.yml")

	if err != nil {
		log.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, &appConfigInstance)

	if err != nil {
		log.Println(err)
	}
}

