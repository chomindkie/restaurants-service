package common

type PlacesSearchResult struct {
	Name    string  `json:"name,omitempty"`
	PlaceID string  `json:"placeId,omitempty"`
	Rating  float32 `json:"rating,omitempty"`
	Image   string  `json:"image,omitempty"`
}
