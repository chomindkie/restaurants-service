package findrestaurants

import "googlemaps.github.io/maps"

type Request struct {
	Keyword string `json:"keyword"`
}

type ResponseModel struct {
	Status Status                     `json:"status"`
	Data   *[]maps.PlacesSearchResult `json:"data,omitempty"`
}

type Status struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
