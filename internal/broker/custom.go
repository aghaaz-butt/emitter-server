package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
)

var redisClient redis.Conn
var client http.Client

type RedisObj struct {
	Key       string
	MembersID []string
}

//InitializeRedisCache connect redis server
func InitializeRedis(s *Service) {

	var err error
	redisClient, err = redis.Dial("tcp", s.Config.RadiurServerIp)
	if err != nil {
		fmt.Println("Redis server is unavailable", err)
	}
	_, err = redisClient.Do("AUTH", s.Config.RadiurServerPassword)
	if err != nil {
		fmt.Println("Redis server is unavailable", err)
	} else {
		fmt.Println("Redis Server Connected!")
	}
}

//Load all Groups from Radius Server
func LoadGroups(s *Service) {

	var GroupsDataJson []map[string]interface{}

	keys, err := redis.Strings(redisClient.Do("KEYS", "*"))
	if err != nil {
		// handle error
	}

	for _, radiusKey := range keys {

		resp, err := redis.String(redisClient.Do("GET", radiusKey))
		if err != nil {
			// handle error
		}

		fmt.Println(resp)
		b := []byte(resp)
		//student := &Student{}
		group := make(map[string]interface{})
		json.Unmarshal(b, group)
		fmt.Println(group)

		var access uint8 = 123

		key, err := s.keygen.CreateKey(s.Config.SecretKey, fmt.Sprintf("%v/", radiusKey), access, time.Unix(9999999999, 0))
		if err != nil {
			fmt.Println("Error ", err)
		}

		GroupDataJson := make(map[string]interface{})
		GroupDataJson["group_id"] = radiusKey
		GroupDataJson["channel_key"] = key

		GroupsDataJson = append(GroupsDataJson, GroupDataJson)

	}

	/*
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
	*/

	BodyParams := make(map[string]interface{})
	BodyParams["auth_token"] = s.Config.AuthKey
	BodyParams["groups_data"] = GroupsDataJson
	jsonValue, _ := json.Marshal(BodyParams)
	fmt.Println(GroupsDataJson)

	req, err := http.NewRequest("POST", s.Config.ApiServer+"v2/UpdateEmitterChannelKeys", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
