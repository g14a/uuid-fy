package pgfunc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"uuid-fy/models"
)

func InsertNeoEventInPG(model models.BlockChainModel) error {
	if CheckIntegrityOfChain() {
		hashJson, err := json.Marshal(model.Data)
		if err != nil {
			log.Println(err)
		}
		
		hash := sha256.Sum256(hashJson)
		
		model.Hash = hex.EncodeToString(hash[:])
		model.TimeStamp = time.Now()
		
		prevHash, err := GetPreviousBlockHash()
		if err != nil {
			return err
		}
		
		if prevHash != "" {
			model.PreviousHash = prevHash
			
			err = pgDB.Insert(&model)
			if err != nil {
				log.Println(err)
			}
		}
		
		return err
	} else {
		return errors.New("the Blockchain is not valid")
	}
}

func GetPreviousBlockHash() (string, error) {
	var previousHash, id string
	err := pgDB.Model(&models.BlockChainModel{}).Column("hash").Column("id").
			  OrderExpr("id DESC").Limit(1).Select(&previousHash, &id)
	
	if previousHash == "" {
		return "genesisblockhash", nil
	}
	
	if err != nil {
		return "", err
	}
	
	return previousHash, nil
}

func CheckIntegrityOfChain() bool {
	var prevHashes, hashes []string
	
	err := pgDB.Model(&models.BlockChainModel{}).
		   Column("previoushash").Select(&prevHashes)
	
	err = pgDB.Model(&models.BlockChainModel{}).
		Column("hash").Select(&hashes)
	
	if err != nil {
		fmt.Println(err)
	}
	
	for k, v := range hashes {
		if k+1 != len(hashes) && v != prevHashes[k+1]{
			return false
		}
	}
	
	return true
}