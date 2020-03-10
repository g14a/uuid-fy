package neofunc

import (
	"encoding/json"
	"fmt"
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

func CreateUser(person models.UserModel) (interface{}, error) {

	defer WriteSession.Close()
	
	person.UUID = uuid.NewV4().String()

	var personInterface map[string]interface{}
	inrec, _ := json.Marshal(person)
	err := json.Unmarshal(inrec, &personInterface)
	if err != nil {
		log.Println(err)
	}

	result, err := WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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

func UpdateContactInfo(name string, person models.ContactInfoModel) (interface{}, error) {

	defer WriteSession.Close()
	
	var personInterface map[string]interface{}
	inrec, _ := json.Marshal(person)
	err := json.Unmarshal(inrec, &personInterface)
	if err != nil {
		log.Println(err)
	}

	result, err := WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH (n:ContactInfoNode { name: $name}) SET n += $props RETURN n", map[string]interface{}{
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

func GetUser(name string) (interface{}, error) {
	defer ReadSession.Close()
	
	result, err := ReadSession.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH(n:UserNode {name: $name}) return n",
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
	defer ReadSession.Close()
	
	result, err := ReadSession.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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

func CreateContactInfo(contactInfo models.ContactInfoModel) (interface{}, error) {
	defer WriteSession.Close()
	
	var contactInterface map[string]interface{}
	bytes, _ := json.Marshal(contactInfo)
	err := json.Unmarshal(bytes, &contactInterface)
	if err != nil {
		log.Println(err)
	}
	
	result, err := WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"CREATE (c:ContactInfoNode { name:$name, phone:$phone, email:$email, address:$address}) RETURN c",
			contactInterface)
		
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

func CreateRelationToContactNode(username, phone string) (interface{}, error) {
	defer WriteSession.Close()
	
	params := map[string]interface{} {
		"username": username,
		"phone": phone,
	}
	
	result, err := WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"Match(a:UserNode),(c:ContactInfoNode) where a.username=$username and c.phone=$phone CREATE (a)-[r:ContactInfoRelation]->(c) return type(r)",
				params)
		
		if err != nil {
			log.Println(err)
		}
		
		if result.Next() {
			fmt.Println(result)
			return result, nil
		}
		
		return nil, result.Err()
	})
	
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	return result, nil
}

func CreateEducationInfo(educationInfo models.EducationInfoModel) (interface{}, error) {
	defer WriteSession.Close()
	
	var educationInterface map[string]interface{}
	bytes, _ := json.Marshal(educationInfo)
	err := json.Unmarshal(bytes, &educationInterface)
	if err != nil {
		log.Println(err)
	}
	
	result, err := WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"CREATE (c:EducationInfoModel { primary:$primary, secondary:$secondary, university:$university}) RETURN c",
			educationInterface)
		
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

