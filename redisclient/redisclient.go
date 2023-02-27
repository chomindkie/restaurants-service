package redisclient

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
	"restaurants-service/common"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

type Cache struct {
	Redis *redis.Client
}

type Cacher interface {
	GetRestaurantsByKeyword(key string) *RestaurantResult
	SaveRestaurantsByKeyword(key string, response []common.PlacesSearchResult, area *maps.LatLng, ttl time.Duration)
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
		logrus.Fatalf("cannot Ping redis: %s", err.Error())
		return nil
	}

	return c
}

func (c Cache) GetRestaurantsByKeyword(key string) *RestaurantResult {
	var restaurantList *RestaurantResult

	key = strings.ToLower(key)
	result, err := c.Redis.HGetAll(key).Result()
	if err != nil {
		return nil
	}

	redisResult := new(RedisResult)
	err = mapstructure.Decode(result, redisResult)
	if err != nil {
		return nil
	}

	if result != nil {
		var restaurants *[]common.PlacesSearchResult
		var area *maps.LatLng

		err = json.Unmarshal([]byte(redisResult.Restaurants), &restaurants)
		if err != nil {
			return nil
		}

		err = json.Unmarshal([]byte(redisResult.Area), &area)
		if err != nil {
			return nil
		}

		restaurantList = &RestaurantResult{
			Restaurants: restaurants,
			Area:        area,
		}
		logrus.Infof("%s Found on Redis", key)
	}
	return restaurantList
}

func (c Cache) SaveRestaurantsByKeyword(key string, response []common.PlacesSearchResult, area *maps.LatLng, ttl time.Duration) {

	key = strings.ToLower(key)
	responseJson, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	areaJson, err := json.Marshal(area)
	if err != nil {
		panic(err)
	}

	c.Redis.HSet(key, "restaurants", responseJson, "area", areaJson)
	logrus.Infof("%s Save on Redis", key)
}
