package education_info

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"uuid-fy/models"
	"uuid-fy/neofunc"
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

