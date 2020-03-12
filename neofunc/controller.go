package neofunc

import (
	"encoding/json"
	"log"
	"uuid-fy/boot"
	"uuid-fy/models"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	uuid "github.com/satori/go.uuid"
)

func CreateUser(person models.UserModel) (interface{}, error) {
	
	person.UUID = uuid.NewV4().String()

	var personInterface map[string]interface{}
	inrec, _ := json.Marshal(person)
	err := json.Unmarshal(inrec, &personInterface)
	if err != nil {
		log.Println(err)
	}

	session, err := boot.GetWriteSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"create(a:UserNode { username:$username, id:$id}) return a",
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


func GetUserUUID(name string) (interface{}, error) {
	session, err := boot.GetReadSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()
	
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH(n:UserNode {username: $username}) return n.id",
			map[string]interface{}{
				"username": name,
			})

		if err != nil {
			log.Println(err)
			return nil, result.Err()
		}

		if result.Next() {
			rmap := result.Record().GetByIndex(0)

			return rmap, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func GetAll() (interface{}, error) {
	session, err := boot.GetReadSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH(n:UserNode) return n",nil)
		
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