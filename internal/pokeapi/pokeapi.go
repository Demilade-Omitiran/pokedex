package pokeapi

import (
	"encoding/json"
	"net/http"
)

type LocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string
		Url  string
	} `json:"results"`
}

func GetLocationAreas(url string) (LocationAreasResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	defer resp.Body.Close()

	var data LocationAreasResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return LocationAreasResponse{}, err
	}

	return data, nil
}
