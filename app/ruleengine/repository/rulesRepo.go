package repository

import (
	"bufio"
	"encoding/json"
	"log"
	"main/app/utils"
	"os"
)

type RulesRepo interface {
	GetAllRules() ([]map[string]interface{}, error)
}

type rulesFromFileRepo struct {
}

var filePath = "data.json"

func NewRulesFromFileRepo() RulesRepo {
	return &rulesFromFileRepo{}
}
func (r rulesFromFileRepo) GetAllRules() ([]map[string]interface{}, error) {
	var err error
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("unable to open file, err: %s", err.Error())
		return nil, utils.ErrUnableToAccessFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var retVal []map[string]interface{}
	for scanner.Scan() {
		var data map[string]interface{}
		b := scanner.Bytes()
		err = json.Unmarshal(b, &data)
		if err != nil {
			log.Println(err)
			continue
		}

		retVal = append(retVal, data)
	}

	return retVal, err
}
