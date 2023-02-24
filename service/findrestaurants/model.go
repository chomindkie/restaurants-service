package findrestaurants

import "restaurants-service/common"

type Request struct {
	Keyword string `json:"keyword"`
}

type ResponseModel struct {
	Status Status                       `json:"status"`
	Data   *[]common.PlacesSearchResult `json:"data,omitempty"`
}

type Status struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
