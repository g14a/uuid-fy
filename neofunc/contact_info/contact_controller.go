package contact_info

import (
	"encoding/json"
	"fmt"
	"log"
	"uuid-fy/boot"
	"uuid-fy/models"
	"uuid-fy/pgfunc"
	
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func CreateContactInfo(contactInfo models.ContactInfoModel) (interface{}, error) {

	session, err := boot.GetWriteSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	var contactInterface map[string]interface{}
	bytes, _ := json.Marshal(contactInfo)
	err = json.Unmarshal(bytes, &contactInterface)
	if err != nil {
		log.Println(err)
	}
	
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"CREATE (c:ContactInfoNode { rootuid: $rootuid, name:$name, phone:$phone, email:$email}) RETURN c",
			contactInterface)

		if err != nil {
			log.Println(err)
		}
		
		var m models.BlockChainModel
		m.Data = models.NeoEvent {
			EventType: "CREATE",
			DataPayload: contactInterface,
			Message: "Add Contact Info",
		}
		
		err = pgfunc.InsertNeoEventInPG(m)
		if err != nil {
			return nil, err
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

	var contactInterface map[string]interface{}
	inrec, _ := json.Marshal(person)
	err := json.Unmarshal(inrec, &contactInterface)
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
			"MATCH (n:ContactInfoNode { name: $name}) SET n += $props RETURN n", map[string]interface{}{
				"name":  name,
				"props": contactInterface,
			})

		if err != nil {
			log.Println(err)
			return nil, result.Err()
		}
		
		var m models.BlockChainModel
		m.Data = models.NeoEvent {
			EventType: "UPDATE",
			DataPayload: contactInterface,
			Message: "Update Contact Info",
		}
		
		err = pgfunc.InsertNeoEventInPG(m)
		if err != nil {
			return nil, err
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

func CreateRelationToContactNode(username, phone string) (interface{}, error) {

	params := map[string]interface{}{
		"username": username,
		"phone":    phone,
	}

	session, err := boot.GetWriteSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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

func DeleteContactInfoOfUser(rootuid string) error {
	params := map[string]interface{}{
		"rootuid": rootuid,
	}

	session, err := boot.GetReadSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		relationDeleteResult, err := tx.Run(
			"MATCH p=(c:ContactInfoNode {rootuid:$rootuid}) DETACH DELETE c;",
			params)

		if err != nil {
			log.Println(err)
		}
		
		var m models.BlockChainModel
		m.Data = models.NeoEvent {
			EventType: "DELETE",
			DataPayload: map[string]interface{} {
				"rootuid": rootuid,
			},
			Message: "Delete Contact Info",
		}
		
		err = pgfunc.InsertNeoEventInPG(m)
		if err != nil {
			return nil, err
		}
		
		return relationDeleteResult, nil
	})

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetContactInfoOfUser(username string) (interface{}, error) {

	params := map[string]interface{}{
		"username": username,
	}

	session, err := boot.GetReadSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"MATCH p=(u:UserNode {username:$username})-[r:ContactInfoRelation]->(c:ContactInfoNode) return c",
			params)

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
