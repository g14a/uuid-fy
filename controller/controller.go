package controller

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	uuid "github.com/satori/go.uuid"
	"log"
	"sync"
	"uuid-fy/boot"
	"uuid-fy/models"
)

var (
	WriteSession neo4j.Session
	ReadSession  neo4j.Session
)

var ReadOnce sync.Once
var WriteOnce sync.Once

func initReadSession()  {
	ReadOnce.Do(func() {
		bootSession, err := boot.GetReadSession()
		if err != nil {
			log.Println(err)
		}
		ReadSession = bootSession
	})
}

func initWriteSession()  {
	WriteOnce.Do(func() {
		bootSession, err := boot.GetWriteSession()
		if err != nil {
			log.Println(err)
		}
		WriteSession = bootSession
	})
}

func init()  {
	initWriteSession()
	initReadSession()
}

func CreatePerson(person models.PersonModel) (interface{}, error) {

	person.UUID = uuid.NewV4().String()

	var personInterface map[string]interface{}
	inrec, _ := json.Marshal(person)
	err := json.Unmarshal(inrec, &personInterface)
	if err != nil {
		log.Println(err)
	}

	result, err := WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"create(a:PersonModel { name:$name, id:$id, dob:$dob, email:$email}) return a",
			personInterface)

		if err != nil {
			log.Println(err)
		}

		if result.Next() {
			rmap := result.Record().GetByIndex(0)

			return rmap.(neo4j.Node).Props(), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func UpdatePerson(name string, person models.UpdatePersonModel) (interface{}, error) {

	var personInterface map[string]interface{}
	inrec, _ := json.Marshal(person)
	err := json.Unmarshal(inrec, &personInterface)
	if err != nil {
		log.Println(err)
	}

	result, err := WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH (n:PersonModel { name: $name}) SET n += $props RETURN n", map[string]interface{}{
				"name": name,
				"props": personInterface,
			})

		if err != nil {
			log.Println(err)
			return nil, result.Err()
		}

		if result.Next() {
			rmap := result.Record().GetByIndex(0)

			return rmap.(neo4j.Node).Props(), nil
		}
		
		return result, nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func GetPerson(name string) (interface{}, error) {
	result, err := ReadSession.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH(n:PersonModel {name: $name}) return n",
			map[string]interface{}{
				"name": name,
			})

		if err != nil {
			log.Println(err)
			return nil, result.Err()
		}

		records, err := neo4j.Collect(result, err)

		var resultMap []map[string]interface{}
		for k, v := range records {
			rmap := v.GetByIndex(k)

			resultMap = append(resultMap, rmap.(neo4j.Node).Props())
		}

		return resultMap, nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func GetAll() (interface{}, error) {
	result, err := ReadSession.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH(n:PersonModel) return n",nil)
		
		if err != nil {
			log.Println(err)
			return nil, result.Err()
		}
		
		records, err := neo4j.Collect(result, err)
		
		var resultMap []map[string]interface{}
		
		for _, v := range records {
			rmap := v.GetByIndex(0)
			
			resultMap = append(resultMap, rmap.(neo4j.Node).Props())
		}
		
		return resultMap, nil
	})
	
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	return result, nil
}