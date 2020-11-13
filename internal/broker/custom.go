package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisObj struct {
	Key       string
	MembersID []string
}

//InitializeRedisCache connect redis server
func InitializeRedis(s *Service) {

	var err error
	s.RedisClient, err = redis.Dial("tcp", s.Config.RadiurServerIp)
	if err != nil {
		fmt.Println("Redis server is unavailable", err)
		return
	}
	_, err = s.RedisClient.Do("AUTH", s.Config.RadiurServerPassword)
	if err != nil {
		fmt.Println("Redis server is unavailable", err)
		return
	}

	fmt.Println("Redis Server Connected!")
	//Load all Groups from Data server
	//LoadGroups(s)

}

//Load all Groups from Radius Server
func LoadGroups(s *Service) {

	var GroupsDataJson []map[string]interface{}

	keys, err := redis.Strings(s.RedisClient.Do("KEYS", "*"))
	if err != nil {
		// handle error
	}

	if len(keys) > 0 {

		for _, radiusKey := range keys {

			resp, err := redis.String(s.RedisClient.Do("GET", radiusKey))
			if err != nil {
				// handle error
			}

			var Group map[string]interface{}
			json.Unmarshal([]byte(resp), &Group)
			fmt.Println(Group)

			var access uint8 = 123
			key, err := s.keygen.CreateKey(s.Config.SecretKey, fmt.Sprintf("%v/", radiusKey), access, time.Unix(9999999999, 0))
			if err != nil {
				// handle error
			}
			fmt.Println("Channel created ", radiusKey)

			if Group["key"] != key {

				GroupDataJson := make(map[string]interface{})
				GroupDataJson["group_id"] = radiusKey
				GroupDataJson["channel_key"] = key

				GroupsDataJson = append(GroupsDataJson, GroupDataJson)
			}
		}

	} else {
		fmt.Println("Nothing found in Radius")
	}

	/*
		BodyParams := make(map[string]interface{})
		BodyParams["auth_token"] = authKey
		jsonValue, _ := json.Marshal(BodyParams)

		req, err := http.NewRequest("POST", apiServer+"v2/GetAllGroups", bytes.NewBuffer(jsonValue))
		req.Heades.Set("Content-Type", "application/json")

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
	*/

	if len(GroupsDataJson) > 0 {

		BodyParams := make(map[string]interface{})
		BodyParams["auth_token"] = s.Config.AuthKey
		BodyParams["groups_data"] = GroupsDataJson
		jsonValue, _ := json.Marshal(BodyParams)

		req, err := http.NewRequest("POST", s.Config.ApiServer+"v2/UpdateEmitterChannelKeys", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.Client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
	}
}

//Authorize User
func AuthorizeUser(Username string, s *Service) bool {

	var UserMap map[string]string
	json.Unmarshal([]byte(Username), &UserMap)
	fmt.Println("---------------------------------------")
	fmt.Println(UserMap)
	fmt.Println("===============")
	fmt.Println(Username)
	fmt.Println("===============")

	resp, err := redis.String(s.RedisClient.Do("GET", "userid_"+UserMap["client_id"]))
	if err != nil {
		fmt.Println("Data not found on Radis")
		return false
	}

	var RadisData map[string]string
	json.Unmarshal([]byte(resp), &RadisData)
	fmt.Println(RadisData)

	fmt.Println("---------------------------------------")

	if UserMap["token"] != RadisData["token"] {
		return false
	}

	return true
}
