package streaming

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type StreamInfo struct {
	Type  string `json:"type"`
	AddOn string `json:"addOn"`
	Price *struct {
		Amount    string `json:"amount"`
		Currency  string `json:"currency"`
		Formatted string `json:"formatted"`
	} `json:"price"`
}

type SearchResult struct {
	Type          string `json:"type"`
	Title         string `json:"title"`
	StreamingInfo struct {
		Ca map[string][]StreamInfo `json:"ca,omitempty"`
	} `json:"streamingInfo"`
	Year         int `json:"year,omitempty"`
	FirstAirYear int `json:"firstAirYear,omitempty"`
}

type SearchResponse struct {
	Result []SearchResult `json:"result"`
}

func Search(title, token string) ([]SearchResult, error) {
	formattedTitle := strings.ReplaceAll(title, " ", "%20")
	url := fmt.Sprintf("https://streaming-availability.p.rapidapi.com/v2/search/title?title=%s&country=ca&output_language=en", formattedTitle)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", token)
	req.Header.Add("X-RapidAPI-Host", "streaming-availability.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := struct {
		Result []SearchResult `json:"result"`
	}{}

	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	return r.Result, nil
}
