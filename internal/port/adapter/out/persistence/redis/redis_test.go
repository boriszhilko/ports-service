package redis

import (
	"context"
	"encoding/json"
	"github.com/boriszhilko/ports-service/internal/port/domain"
	"os"
	"reflect"
	"testing"
)

var redisRepository = NewRedisRepository(os.Getenv("REDIS_URL"))
var testPort = domain.Port{
	ID:          "123",
	Name:        "test1",
	City:        "test2",
	Country:     "test3",
	Alias:       []string{"test4"},
	Regions:     []string{"test5"},
	Coordinates: []float64{6},
	Province:    "test7",
	Timezone:    "test8",
	Unlocs:      []string{"test9"},
	Code:        "test10",
}

func TestNewRedisRepository(t *testing.T) {
	_, err := redisRepository.client.Ping(context.Background()).Result()
	if err != nil {
		t.Errorf("Error pinging redis: %v", err)
	}
}

func TestRedisRepository_Create(t *testing.T) {
	// WHEN
	err := redisRepository.CreateOrUpdate(context.TODO(), []domain.Port{testPort})
	if err != nil {
		t.Errorf("Error creating port: %v", err)
	}
	// THEN
	actualPort := get(t, testPort.ID)
	if reflect.DeepEqual(testPort, actualPort) == false {
		t.Errorf("Expected port: %v, actual port: %v", testPort, actualPort)
	}
}

func TestRedisRepository_Update(t *testing.T) {
	// GIVEN
	updatedPort := domain.Port{
		ID:   "123",
		Name: "updatedName",
	}
	add(t, testPort)
	// WHEN
	err := redisRepository.CreateOrUpdate(context.TODO(), []domain.Port{updatedPort})
	if err != nil {
		t.Errorf("Error updating port: %v", err)
	}
	// THEN
	actualPort := get(t, updatedPort.ID)
	if actualPort.Name != updatedPort.Name {
		t.Errorf("Expected name: %v, actual name: %v", updatedPort.Name, actualPort.Name)
	}
}

func get(t *testing.T, id string) domain.Port {
	portJson, err := redisRepository.client.Get(context.Background(), id).Result()
	if err != nil {
		t.Errorf("Error getting port: %v", err)
	}
	var actualPort domain.Port
	err = json.Unmarshal([]byte(portJson), &actualPort)
	if err != nil {
		t.Errorf("Error unmarshalling port: %v", err)
	}
	return actualPort
}

func add(t *testing.T, port domain.Port) {
	portJson, err := json.Marshal(port)
	if err != nil {
		t.Errorf("Error marshalling port: %v", err)
	}
	err = redisRepository.client.Set(context.Background(), port.ID, portJson, 0).Err()
	if err != nil {
		t.Errorf("Error setting port: %v", err)
	}
}
