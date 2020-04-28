package health_info

import (
	"encoding/json"
	"log"
	"uuid-fy/boot"
	"uuid-fy/models"
	"uuid-fy/pgfunc"
	
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func CreateHealthInfo(healthInfo models.HealthInfoModel) (interface{}, error) {

	session, err := boot.GetWriteSession()
	if err != nil {
		log.Print(err)
	}

	defer session.Close()

	var healthInterface map[string]interface{}
	bytes, _ := json.Marshal(healthInfo)
	err = json.Unmarshal(bytes, &healthInterface)
	if err != nil {
		log.Println(err)
	}

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"CREATE (h:HealthInfoNode { rootuid:$rootuid, birthhospital:$birthhospital, bloodgroup:$bloodgroup}) RETURN h",
			healthInterface)

		if err != nil {
			log.Println(err)
		}
		
		var m models.BlockChainModel
		m.Data = models.NeoEvent {
			EventType: "CREATE",
			DataPayload: healthInterface,
			Message: "Add Health Info",
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

func CreateRelationToHealthNode(uuid string) (interface{}, error) {
	
	params := map[string]interface{}{
		"id": uuid,
	}
	
	session, err := boot.GetWriteSession()
	if err != nil {
		log.Print(err)
	}
	
	defer session.Close()
	
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		result, err := tx.Run(
			"Match(a:UserNode),(h:HealthInfoNode) where a.id=$id and h.rootuid=$id CREATE (a)-[r:HealthInfoRelation]->(h) return type(r)",
			params)
		
		if err != nil {
			log.Println(err)
		}
		
		if result.Next() {
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

func GetHealthInfoOfUser(username string) (interface{}, error) {
	
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
			"MATCH p=(u:UserNode {username:$username})-[r:HealthInfoRelation]->(h:HealthInfoNode) return h;",
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
