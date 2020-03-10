package education_info

import (
	"encoding/json"
	"fmt"
	"log"
	"uuid-fy/models"
	"uuid-fy/neofunc"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func CreateEducationInfo(educationInfo models.EducationInfoModel) (interface{}, error) {
	defer neofunc.WriteSession.Close()

	var educationInterface map[string]interface{}
	bytes, _ := json.Marshal(educationInfo)
	err := json.Unmarshal(bytes, &educationInterface)
	if err != nil {
		log.Println(err)
	}

	result, err := neofunc.WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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
	defer neofunc.WriteSession.Close()

	params := map[string]interface{}{
		"id": uuid,
	}

	result, err := neofunc.WriteSession.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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
	defer neofunc.ReadSession.Close()

	params := map[string]interface{}{
		"username": username,
	}

	result, err := neofunc.ReadSession.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
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
