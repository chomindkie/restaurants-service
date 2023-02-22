package findrestaurants

import (
	"context"
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

type Service struct {
	cache redisclient.Cacher
}

func NewService(cache redisclient.Cacher) *Service {
	return &Service{
		cache: cache,
	}
}

func (s *Service) FindRestaurant(c echo.Context, request Request) (*maps.PlacesSearchResponse, error) {
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
		lat, lng := findLatLngByKeyword(keyword)

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
			logrus.Errorf("Error making request: %v", err)
			return nil, errs.JSON(c, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, ""))
		}
		results = &res

		// Save to Redis
		ttl, err := time.ParseDuration(viper.GetString("redis.ttl"))
		if err != nil {
			log.Errorf("ParseDuration redis.redeeming-ttl error: %s", err.Error())
			return nil, errs.NewStatus(http.StatusInternalServerError, errs.INTERNAL_ERROR)
		}
		s.cache.SaveRestaurantsByKeyword(keyword, res, ttl)
	}

	return results, nil
}

func findLatLngByKeyword(keyword string) (float64, float64) {
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
		logrus.Errorf("Error making request: %v", err)
		return 13.809082, 100.537801 // Default to Bang Sue
	}

	// Get the latitude and longitude
	lat := resp[0].Geometry.Location.Lat
	lng := resp[0].Geometry.Location.Lng

	logrus.Printf("Find Restaurant in area %s Latitude: %f Longitude: %f", keyword, lat, lng)

	return lat, lng
}
