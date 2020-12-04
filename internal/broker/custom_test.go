// https://itnext.io/golang-testing-mocking-redis-b48d09386c70

package broker

import (
	"github.com/elliotchance/redismock"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

var (
	key = "key"
	val = "val"
)

//InitializeRedisCache connect redis server
func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	code := m.Run()
	os.Exit(code)
}

//Set dummy value in redis for time.Duration(0)
func TestInitializeRedis(t *testing.T) {
	exp := time.Duration(0)

	mock := redismock.NewNiceMock(client)
	r := NewRedisRepository(mock)

	mock.On("Set", key, val, exp).Return(redis.NewStatusResult("", nil))
	err := r.Set(key, val, exp)
	assert.NoError(t, err)

}

// Get value against assigned key
func TesAuthorizeUser(t *testing.T) {
	mock := redismock.NewNiceMock(client)
	r := NewRedisRepository(mock)

	mock.On("Get", key).Return(redis.NewStringResult(val, nil))
	res, err := r.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, res)
}