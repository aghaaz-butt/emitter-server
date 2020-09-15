package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/emitter-io/emitter/internal/config"
	"github.com/emitter-io/emitter/internal/service/keygen"
)

type Service struct {
	Config        *config.Config     // The configuration for the service.
	keygen        *keygen.Service    // The key generation provider.
}


//Load all Groups from Data server
func LoadGroups(s *Service) {

	apiServer	:= s.Config.ApiServer
	secretKey	:= s.Config.SecretKey
	authKey 	:= s.Config.AuthKey

	BodyParams := make(map[string]interface{})
	BodyParams["auth_token"] = authKey
	jsonValue, _ := json.Marshal(BodyParams)

	req, err := http.NewRequest("POST", apiServer+"v2/GetAllGroups", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if resp.Status != "200 OK" {
		panic(err)
	}

	var GroupsData map[string]interface{}
	if err := json.Unmarshal(body, &GroupsData); err != nil {
		panic(err)
	}

	if fmt.Sprintf("%v", GroupsData["status"]) != "200" {
		fmt.Println(GroupsData["status"])
	}

	Groups := GroupsData["groups"].([]interface{})

	var GroupsDataJson []map[string]interface{}

	for i := 0; i < len(Groups); i++ {

		Group := Groups[i].(map[string]interface{})

		var access uint8 = 123
		key := ""

		key, err := s.keygen.CreateKey(secretKey, fmt.Sprintf("%v/", Group["id"]), access, time.Unix(9999999999, 0))
		if err != nil {
			fmt.Println(err)
		}

		GroupDataJson := make(map[string]interface{})
		GroupDataJson["group_id"] = Group["id"]
		GroupDataJson["channel_key"] = key

		GroupsDataJson = append(GroupsDataJson, GroupDataJson)

	}

	BodyParams = make(map[string]interface{})
	BodyParams["auth_token"] = authKey
	BodyParams["groups_data"] = GroupsDataJson
	jsonValue, _ = json.Marshal(BodyParams)

	req, err = http.NewRequest("POST", apiServer+"v2/UpdateEmitterChannelKeys", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
