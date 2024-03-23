package client

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type YouTubeAPIClient struct {
	APIKeys      []string
	CurrentIndex int
}

func NewYouTubeAPIClient(apiKeys []string) *YouTubeAPIClient {
	return &YouTubeAPIClient{
		APIKeys:      apiKeys,
		CurrentIndex: 0,
	}
}

func (c *YouTubeAPIClient) GetNextAPIKey() string {
	key := c.APIKeys[c.CurrentIndex]
	c.CurrentIndex = (c.CurrentIndex + 1) % len(c.APIKeys)
	return key
}

// YTClient to fetch videos.
func (c *YouTubeAPIClient) GetYoutubeVideoDetails(query, pageToken string) (map[string]interface{}, error) {
	/*
		Handle rate limiting scenario,
		it distributes call over the keys in round robin fashion.
	*/
	apiKey := c.GetNextAPIKey()

	searchQuery := query

	baseURL := "https://www.googleapis.com/youtube/v3/search"
	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("maxResults", "5")
	params.Set("q", searchQuery)
	params.Set("key", apiKey)
	params.Set("order", "date")
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}

	url := baseURL + "?" + params.Encode()

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching data from YouTube API:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Error decoding JSON response:", err)
		return nil, err
	}
	return response, nil
}
