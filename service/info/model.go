package info

type InfoResponse struct {
	Build Build `json:"build"`
}

type Build struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}