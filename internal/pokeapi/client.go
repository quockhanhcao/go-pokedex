package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
