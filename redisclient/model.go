package redisclient

import (
	"googlemaps.github.io/maps"
	"restaurants-service/common"
)

type RedisResult struct {
	Restaurants string `json:"restaurants" form:"restaurants"`
	Area        string `json:"area" form:"area"`
}

type RestaurantResult struct {
	Restaurants *[]common.PlacesSearchResult `json:"restaurants" form:"restaurants"`
	Area        *maps.LatLng                 `json:"area" form:"area"`
}
