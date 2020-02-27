package main

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"uuid-fy/boot"
)

func main() {

}

func helloWorld() (string, error) {

	 session, err := boot.GetSession()
	 if err != nil {
		log.Println(err)
	 }

	 greeting, err := session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
	 	result, err := tx.Run(
			"CREATE (a:Greeting) SET a.message=$message RETURN a.message",
			map[string]interface{}{"message": "hello, prabhakar"})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}
		return nil, result.Err()
	 })

	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}
