package common

type PlacesSearchResult struct {
	Name     string  `json:"name,omitempty"`
	PlaceID  string  `json:"placeId,omitempty"`
	Rating   float32 `json:"rating,omitempty"`
	Image    string  `json:"image,omitempty"`
	Location LatLng  `json:"area,omitempty"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
