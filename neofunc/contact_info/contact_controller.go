package contact_info

import (
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"uuid-fy/models"
	"uuid-fy/neofunc"
)

func CreateContactInfo(contactInfo models.ContactInfoModel) (interface{}, error) {
	defer neofunc.WriteSession.Close()
	
	var contactInterface map[string]interface{}
	bytes, _ := json.Marshal(contactInfo)
	err := json.Unmarshal(bytes, &contactInterface)
	if err != nil {
		log.Println(err)
	}
	
	result, err := neofunc.WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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

func UpdateContactInfo(name string, person models.ContactInfoModel) (interface{}, error) {
	
	defer neofunc.WriteSession.Close()
	
	var personInterface map[string]interface{}
	inrec, _ := json.Marshal(person)
	err := json.Unmarshal(inrec, &personInterface)
	if err != nil {
		log.Println(err)
	}
	
	result, err := neofunc.WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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

func CreateRelationToContactNode(username, phone string) (interface{}, error) {
	defer neofunc.WriteSession.Close()
	
	params := map[string]interface{} {
		"username": username,
		"phone": phone,
	}
	
	result, err := neofunc.WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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

func GetContactInfoOfUser(username string) (interface{}, error) {
	defer neofunc.ReadSession.Close()
	
	params := map[string]interface{} {
		"username": username,
	}
	
	result, err := neofunc.ReadSession.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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