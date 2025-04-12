package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func GetLocationAreaData(url string) (Response, error) {
	data := Response{}
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
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
		return Response{}, err
	}
	return data, nil
}
