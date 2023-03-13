package models

import (
	"encoding/json"
	"log"
	"os"
)

type Cluster struct {
	Nickname  string  `json:"nickname"`
	Domain    string  `json:"domain"`
	ApiKey    string  `json:"api_key"`
	ClusterId string  `json:"cluster_id"`
	Id        *string `json:"id,omitempty"`
}

var ClusterMap = map[string]Cluster{}

func init() {

	//path, err := os.Getwd()
	//if err != nil {
	//	log.Println(err)
	//	os.Exit(2)
	//}
	//configFilePath := path + "/app/cluster/models/del.txt"

	configFilePath := "/opt/cluster.conf"

	jsonFile, err := os.Open(configFilePath)

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	defer jsonFile.Close()

	err = json.NewDecoder(jsonFile).Decode(&ClusterMap)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	return
}
