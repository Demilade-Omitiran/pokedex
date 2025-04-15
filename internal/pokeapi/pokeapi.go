package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var cache *Cache

func getCache() *Cache {
	if cache == nil {
		cache = NewCache(1 * time.Minute)
	}

	return cache
}

func GetLocationAreas(url string) (LocationAreasResponse, error) {
	apiCache := getCache()
	cachedData, found := apiCache.Get(url)

	if found {
		fmt.Println("Cache hit for key:", url)
		var data LocationAreasResponse
		if err := json.Unmarshal(cachedData, &data); err != nil {
			return LocationAreasResponse{}, err
		}
		return data, nil
	}

	fmt.Println("Cache miss for key:", url)

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	defer resp.Body.Close()

	byteData, err := io.ReadAll(resp.Body)

	if err != nil {
		return LocationAreasResponse{}, err
	}

	var data LocationAreasResponse
	if err := json.Unmarshal(byteData, &data); err != nil {
		return LocationAreasResponse{}, err
	}

	apiCache.Add(url, byteData)

	return data, nil
}

func GetPokemonInLocationArea(name string) (PokemonInLocationAreaResponse, error) {
	fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
	apiCache := getCache()

	cachedData, found := apiCache.Get(fullUrl)

	if found {
		fmt.Println("Cache hit for key:", fullUrl)
		var data PokemonInLocationAreaResponse
		if err := json.Unmarshal(cachedData, &data); err != nil {
			return PokemonInLocationAreaResponse{}, err
		}
		return data, nil
	}

	fmt.Println("Cache miss for key:", fullUrl)

	resp, err := http.Get(fullUrl)

	if err != nil {
		return PokemonInLocationAreaResponse{}, err
	}

	defer resp.Body.Close()

	byteData, err := io.ReadAll(resp.Body)

	if err != nil {
		return PokemonInLocationAreaResponse{}, err
	}

	var data PokemonInLocationAreaResponse
	if err := json.Unmarshal(byteData, &data); err != nil {
		return PokemonInLocationAreaResponse{}, err
	}

	apiCache.Add(fullUrl, byteData)
	return data, nil
}

func GetPokemonInfo(name string) (Pokemon, error) {
	fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	apiCache := getCache()

	cachedData, found := apiCache.Get(fullUrl)

	if found {
		fmt.Println("Cache hit for key:", fullUrl)
		var data Pokemon
		if err := json.Unmarshal(cachedData, &data); err != nil {
			return Pokemon{}, err
		}
		return data, nil
	}

	fmt.Println("Cache miss for key:", fullUrl)

	resp, err := http.Get(fullUrl)

	if err != nil {
		return Pokemon{}, err
	}

	defer resp.Body.Close()

	byteData, err := io.ReadAll(resp.Body)

	if err != nil {
		return Pokemon{}, err
	}

	var data Pokemon

	if err := json.Unmarshal(byteData, &data); err != nil {
		return Pokemon{}, err
	}

	apiCache.Add(fullUrl, byteData)
	return data, nil
}
