package findrestaurants

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"googlemaps.github.io/maps"
	"net/http"
	"restaurants-service/library/errs"
	"restaurants-service/redisclient"
	"time"
)

var (
	SUCCESS = Status{Code: "SUCCESS", Message: "success"}
)

type Service struct {
	cache redisclient.Cacher
}

func NewService(cache redisclient.Cacher) *Service {
	return &Service{
		cache: cache,
	}
}

func (s *Service) FindRestaurant(c echo.Context, request Request) (*ResponseModel, error) {
	// Specify the search criteria
	radius := 1000 // in meters
	types := "restaurant"
	keywordDefault := "Bang Sue"

	keyword := request.Keyword
	var results *maps.PlacesSearchResponse

	if keyword == "" {
		keyword = keywordDefault
	}

	// Find in redis
	results = s.cache.GetRestaurantsByKeyword(keyword)

	if results == nil {
		apiKey := viper.GetString("apiKey")
		client, err := maps.NewClient(maps.WithAPIKey(apiKey))
		if err != nil {
			panic(err)
		}

		// Find location
		lat, lng, err := findLatLngByKeyword(keyword)

		if err != nil {
			return nil, err
		}

		// Set the location to search around
		location := &maps.LatLng{Lat: lat, Lng: lng} // Default location: Bang Sue

		// Call the nearby search
		res, err := client.NearbySearch(context.Background(), &maps.NearbySearchRequest{
			Location: location,
			Radius:   uint(radius),
			Type:     maps.PlaceType(types),
			Keyword:  keyword,
		})

		if err != nil {
			logrus.Errorf("Error making request to NearbySearch: %v", err)
			return nil, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, fmt.Sprintf("Error making request to NearbySearch: %v", err))
		}
		results = &res

		// Save to Redis
		ttl, err := time.ParseDuration(viper.GetString("redis.ttl"))
		if err != nil {
			log.Errorf("ParseDuration redis.redeeming-ttl error: %s", err)
			return nil, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, fmt.Sprintf("ParseDuration redis.redeeming-ttl error: %s", err))
		}
		s.cache.SaveRestaurantsByKeyword(keyword, res, ttl)
	}

	res := &ResponseModel{
		Status: SUCCESS,
		Data:   results,
	}
	return res, nil
}

func findLatLngByKeyword(keyword string) (float64, float64, error) {
	apiKey := viper.GetString("apiKey")

	// Create client
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		logrus.Fatalf("Error creating client: %v", err)
	}

	// Create a request to the Geocoding API
	r := &maps.GeocodingRequest{
		Address: keyword,
	}

	// Call the Geocode function
	resp, err := client.Geocode(context.Background(), r)
	if err != nil {
		logrus.Errorf("Error making request to Geocode: %v", err)
		return 0, 0, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, fmt.Sprintf("Error making request to Geocode: %v", err))
	}

	if len(resp) == 0 {
		logrus.Errorf("Not found Latitude Longitude of : %s", keyword)
		return 0, 0, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, fmt.Sprintf("Not found Latitude Longitude of : %s", keyword))
	}

	// Get the latitude and longitude
	lat := resp[0].Geometry.Location.Lat
	lng := resp[0].Geometry.Location.Lng

	logrus.Printf("Find Restaurant in area %s Latitude: %f Longitude: %f", keyword, lat, lng)

	return lat, lng, nil
}
