package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Name struct {
	Short string
}

type Team struct {
	Names Name
}

type GameResponse struct {
	Home              Team
	Away              Team
	FinalMessage      string
	BracketPositionId string
}

type Response struct {
	Games []*GameResponse
}

func parseBracket(response []byte) *Response {
	res := Response{}

	if err := json.Unmarshal(response, &res); err != nil {
		panic(err)
	}
	return &res
}

func getUrl(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func Scrape(url string) ([]*GameResponse, error) {
	responseBody, err := getUrl(url)
	if err != nil {
		return nil, err
	}
	response := parseBracket(responseBody)
	return response.Games, nil
}
