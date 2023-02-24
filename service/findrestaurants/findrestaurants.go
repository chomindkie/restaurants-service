package findrestaurants

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"googlemaps.github.io/maps"
	"net/http"
	"restaurants-service/common"
	"restaurants-service/library/errs"
	"restaurants-service/redisclient"
	"strings"
	"time"
)

var (
	SUCCESS = Status{Code: "SUCCESS", Message: "success"}
)

type Service struct {
	place *maps.Client
	cache redisclient.Cacher
}

func NewService(place *maps.Client, cache redisclient.Cacher) *Service {
	return &Service{
		place: place,
		cache: cache,
	}
}

func (s *Service) GetListOfRestaurantByKeyword(request Request) (*ResponseModel, error) {
	// Specify the search criteria
	radius := 1000 // in meters
	types := "restaurant"
	keywordDefault := "Bang Sue"

	keyword := request.Keyword
	var results *[]common.PlacesSearchResult

	if keyword == "" {
		keyword = keywordDefault
	}

	// Find in redis
	results = s.cache.GetRestaurantsByKeyword(keyword)

	if results == nil {

		// Find location
		lat, lng, err := findLatLngByKeyword(s.place, keyword)
		if err != nil {
			return nil, err
		}

		// Set the location to search around
		location := &maps.LatLng{Lat: lat, Lng: lng} // Default location: Bang Sue

		req := maps.NearbySearchRequest{
			Location: location,
			Radius:   uint(radius),
			Type:     maps.PlaceType(types),
			Keyword:  keyword,
		}

		// Call the nearby search
		res, err := s.place.NearbySearch(context.Background(), &req)
		if err != nil {
			logrus.Errorf("Error making request to NearbySearch: %v", err)
			return nil, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, fmt.Sprintf("Error making request to NearbySearch: %v", err))
		}

		// fist list of restaurants"s result
		searchResultDto := getPlacesSearchResultDto(res)

		// Check if there are additional pages of results
		if len(res.NextPageToken) > 0 {
			var errors error
			allSearchResultDto, errors := findRestaurants(s.place, req, res.NextPageToken, searchResultDto, 0)
			if errors == nil {
				searchResultDto = allSearchResultDto
			}
		}

		logrus.Infof("Total %v found with keyword: %s", len(searchResultDto), keyword)
		results = &searchResultDto

		// Save to Redis
		ttl, err := time.ParseDuration(viper.GetString("redis.ttl"))
		if err != nil {
			log.Errorf("ParseDuration redis.redeeming-ttl error: %s", err)
			return nil, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, fmt.Sprintf("ParseDuration redis.redeeming-ttl error: %s", err))
		}
		s.cache.SaveRestaurantsByKeyword(keyword, searchResultDto, ttl)
	}

	res := &ResponseModel{
		Status: SUCCESS,
		Data:   results,
	}
	return res, nil
}

func findRestaurants(client *maps.Client, request maps.NearbySearchRequest, nextPageToken string, list []common.PlacesSearchResult, attempts int) ([]common.PlacesSearchResult, error) {
	var nextPageResp maps.PlacesSearchResponse
	var nextPageRequest maps.NearbySearchRequest
	maxAttempt := viper.GetInt("max-attempt")

	if len(nextPageToken) == 0 {
		return list, nil
	} else {
		maskToken := "***" + nextPageToken[len(nextPageToken)-3:]
		// Set up the search request for the next page of results
		nextPageRequest := maps.NearbySearchRequest{
			PageToken: nextPageToken,
			Location:  request.Location,
			Radius:    request.Radius,
			Type:      request.Type,
			Keyword:   request.Keyword,
		}

		/*
		   Due to issue in google api that
		   there is a short delay between when a next_page_token is issued, and when it will become valid,
		   so we need to retry until it success
		*/
		time.Sleep(100 * time.Millisecond)
		var err error
		nextPageResp, err = client.NearbySearch(context.Background(), &nextPageRequest)

		if err != nil {
			if attempts == maxAttempt {
				logrus.Errorf("Error making request to NearbySearch: %v %s with token %s\n", err.Error(), request.Keyword, maskToken)
				return nil, errs.New(http.StatusInternalServerError, errs.INTERNAL_ERROR.Code, fmt.Sprintf("Error making request to NearbySearch: %v", err))
			} else if attempts != 20 {
				attempts += 1
				logrus.Errorf("Retry %v with pageToken %v\n", attempts, maskToken)
				list, err = findRestaurants(client, request, nextPageToken, list, attempts)
			}
		} else {
			logrus.Infof("Get list of restaurant success with pageToken %v\n", maskToken)
			getPlacesSearchResultDto(nextPageResp)
			list = append(list, getPlacesSearchResultDto(nextPageResp)...)
		}

	}
	return findRestaurants(client, nextPageRequest, nextPageResp.NextPageToken, list, 0)
}

func findLatLngByKeyword(client *maps.Client, keyword string) (float64, float64, error) {

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

func getPlacesSearchResultDto(response maps.PlacesSearchResponse) []common.PlacesSearchResult {
	var list []common.PlacesSearchResult
	for _, item := range response.Results {
		placesSearchResult := common.PlacesSearchResult{
			Name:    item.Name,
			PlaceID: item.PlaceID,
			Rating:  item.Rating,
			Image:   getImageUrl(item.Photos, item.Icon),
		}
		list = append(list, placesSearchResult)
	}
	return list
}

func getImageUrl(photo []maps.Photo, icon string) string {
	if len(photo) != 0 {
		return strings.Replace(viper.GetString("imageUrl"), "{photoReference}", photo[0].PhotoReference, -1)
	}
	return icon
}
