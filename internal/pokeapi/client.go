package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/quockhanhcao/go-pokedex/internal/pokecache"
)

type Response struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_details"`
		}
	}
	GameIndex int `json:"game_index"`
	Id        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	}
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			}
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Weight         int    `json:"weight"`
	Height         int    `json:"height"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

var cache = pokecache.NewCache(10 * time.Second)

func GetLocationAreaData(url string) (Response, error) {
	data := Response{}
	cacheData, ok := cache.Get(url)
	if ok {
		fmt.Println("Cache hit")
		err := json.Unmarshal(cacheData, &data)
		if err != nil {
			return Response{}, err
		}
		return data, nil
	}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return Response{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
		return Response{}, err
	}
	// add to cache
	cache.Add(url, body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
		return Response{}, err
	}
	return data, nil
}

func GetLocationDetailedData(locationName string) (LocationResponse, error) {
	data := LocationResponse{}
	cacheData, ok := cache.Get(locationName)
	if ok {
		err := json.Unmarshal(cacheData, &data)
		if err != nil {
			return LocationResponse{}, err
		}
		return data, nil
	}
	var sb strings.Builder
	sb.WriteString("https://pokeapi.co/api/v2/location-area/")
	sb.WriteString(locationName)
	url := sb.String()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return LocationResponse{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
		return LocationResponse{}, err
	}
	cache.Add(locationName, body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
		return LocationResponse{}, err
	}
	return data, nil
}

func GetPokemonStats(pokemon string) (Pokemon, error) {
	data := Pokemon{}
	cacheData, ok := cache.Get(pokemon)
	if ok {
		err := json.Unmarshal(cacheData, &data)
		if err != nil {
			return Pokemon{}, err
		}
		return data, nil
	}
	var sb strings.Builder
	sb.WriteString("https://pokeapi.co/api/v2/pokemon/")
	sb.WriteString(pokemon)
	url := sb.String()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return Pokemon{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
		return Pokemon{}, err
	}
	cache.Add(pokemon, body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
		return Pokemon{}, err
	}
	return data, nil
}

func CloseCache() {
	cache.Close()
}
