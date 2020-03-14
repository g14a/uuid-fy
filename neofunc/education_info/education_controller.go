package education_info

import (
	"encoding/json"
	"fmt"
	"log"
	"uuid-fy/boot"
	"uuid-fy/models"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func CreateEducationInfo(educationInfo models.EducationInfoModel) (interface{}, error) {

	var educationInterface map[string]interface{}
	bytes, _ := json.Marshal(educationInfo)
	err := json.Unmarshal(bytes, &educationInterface)
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
			"CREATE (e:EducationInfoNode { rootuid:$rootuid, primary:$primary, secondary:$secondary, university:$university}) RETURN e",
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

func CreateRelationToEducationNode(uuid string) (interface{}, error) {

	params := map[string]interface{}{
		"id": uuid,
	}

	session, err := boot.GetWriteSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	fmt.Println(params);
	
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"Match(a:UserNode),(e:EducationInfoNode) where a.id=$id and e.rootuid=$id CREATE (a)-[r:EducationInfoRelation]->(e) return type(r)",
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

func GetEducationInfoOfUser(username string) (interface{}, error) {

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
			"MATCH p=(u:UserNode {username:$username})-[r:EducationInfoRelation]->(e:EducationInfoNode) return e;",
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
