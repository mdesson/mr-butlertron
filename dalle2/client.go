package dalle2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	imageURL    = "https://api.openai.com/v1/images/generations"
	imagesCount = 9
	imageSize   = "1024x1024"
)

type Client struct {
	token      string
	httpClient http.Client
}

func NewClient(token string) *Client {
	return &Client{token: token, httpClient: http.Client{}}
}

type ImageRequest struct {
	Prompt string `json:"prompt"`
	Number int    `json:"n"`
	Size   string `json:"size"`
}

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func (c *Client) RequestImages(prompt string) ([]string, error) {
	imageReq := ImageRequest{
		Prompt: prompt,
		Number: imagesCount,
		Size:   imageSize,
	}

	reqBody, err := json.Marshal(imageReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, imageURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyStr, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s\n", resp.StatusCode, bodyStr)
	}

	defer resp.Body.Close()

	var imageRes ImageResponse
	if err = json.NewDecoder(resp.Body).Decode(&imageRes); err != nil {
		return nil, err
	}
	
	URLs := make([]string, 0)
	for _, data := range imageRes.Data {
		URLs = append(URLs, data.Url)
	}

	return URLs, nil
}
