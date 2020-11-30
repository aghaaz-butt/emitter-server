package broker

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	cfg "github.com/emitter-io/emitter/internal/config"

)

var redisClient redis.Conn


//InitializeRedisCache connect redis server
func TestInitializeRedis(t *testing.T) {

	var err error
	c := &cfg.Config{RadiurServerIp: "51.68.162.55:6379",
		RadiurServerPassword: "Mch%5$$3)83Kdh1!a",}

	// checking redis connection
	redisClient, err = redis.Dial("tcp", c.RadiurServerIp)
	if err != nil {
		fmt.Println("Redis server is unavailable", err)
		assert.Nil(t, err)
	}
	assert.Nil(t, err)

	// checking redisClient with password
	_, err = redisClient.Do("AUTH", c.RadiurServerPassword)
	if err != nil {
		fmt.Println("Password of redis is incorrect", err)
		assert.Nil(t, err)
	}
	assert.Nil(t, err)

}

//InitializeRedisCache connect redis server
func TestAuthorizeUser(t *testing.T) {

	// generating random number string
	randonNumber := strconv.Itoa(rand.Intn(10000))

	// updating random number string in redis
	redisClient.Do("SET",  randonNumber + "_user_id_" + randonNumber, randonNumber)

	// getting random number string from redis
	resp, err := redis.String(redisClient.Do("GET", randonNumber + "_user_id_" + randonNumber))
	if err != nil {
		fmt.Println("user isn't available in redis", err)
		assert.Nil(t, err)
	}

	// checking random number string
	assert.Equal(t, randonNumber, resp, "user is available in redis")

	// remove random value from redis
	redisClient.Do("SET",  randonNumber + "_user_id_" + randonNumber, "")
}
