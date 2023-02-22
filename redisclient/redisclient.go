package redisclient

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
	"time"

	"github.com/go-redis/redis/v7"
)

type Cache struct {
	Redis *redis.Client
}

type Cacher interface {
	GetRestaurantsByKeyword(key string) *[]maps.PlacesSearchResult
	SaveRestaurantsByKeyword(key string, response []maps.PlacesSearchResult, ttl time.Duration)
}

func New(redis *redis.Client) *Cache {
	return &Cache{
		Redis: redis,
	}
}

func NewCache(address string, dbIndex int) *Cache {
	c := &Cache{}
	options := redis.Options{
		Addr: address,
		DB:   dbIndex,
	}

	c.Redis = redis.NewClient(&options)
	_, err := c.Redis.Ping().Result()
	if err != nil {
		//acnlog.Fatalf("cannot Ping redis: %s", err.Error())
		return nil
	}

	return c
}

func (c Cache) GetRestaurantsByKeyword(key string) *[]maps.PlacesSearchResult {
	// Create a new person object
	var response []maps.PlacesSearchResult

	result, err := c.Redis.Get(key).Result()
	if err != nil {
		return nil
	}

	// Unmarshal the JSON string into the PlacesSearchResponse
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		return nil
	}

	logrus.Infof("%s Found on Redis", key)
	return &response
}

func (c Cache) SaveRestaurantsByKeyword(key string, response []maps.PlacesSearchResult, ttl time.Duration) {
	// Convert the PlacesSearchResponse to JSON
	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	c.Redis.Set(key, json, ttl)
	logrus.Infof("%s Save on Redis", key)
}
