package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"personal-project/fampay-assignment/models"
	"personal-project/fampay-assignment/repo"
	"sync"
	"time"
)

type ServiceImpl struct {
	repo repo.IRepository
}

func NewService(repo repo.IRepository) IService {
	return &ServiceImpl{
		repo: repo,
	}
}

func fetchLatestVideos(query, pageToken string, resultChan chan<- []models.Video) (string, error) {

	apiKey := os.Getenv("YTAPIKEY")
	log.Println("apiKey", apiKey)

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
		return "", err
	}
	defer resp.Body.Close()

	log.Println("Response", resp)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Error decoding JSON response:", err)
		return "", err
	}
	var nextPageToken string

	videos := make([]models.Video, 0)
	// if response["items"] == nil {
	// 	return response["nextPageToken"].(string), nil
	// }
	for _, item := range response["items"].([]interface{}) {
		videoData := item.(map[string]interface{})["snippet"].(map[string]interface{})
		title := videoData["title"].(string)
		description := videoData["description"].(string)
		publishedAt, _ := time.Parse(time.RFC3339, videoData["publishedAt"].(string))
		videoId := item.(map[string]interface{})["id"].(map[string]interface{})["videoId"].(string)
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)
		nextPageToken = response["nextPageToken"].(string)
		thumbnails := make([]string, 0)
		for _, thumb := range videoData["thumbnails"].(map[string]interface{}) {
			thumbnails = append(thumbnails, thumb.(map[string]interface{})["url"].(string))
		}
		videos = append(videos, models.Video{
			Id:          videoId,
			Title:       title,
			Description: description,
			PublishedAt: publishedAt,
			Thumbnails:  thumbnails,
			URL:         videoURL,
		})
	}

	resultChan <- videos
	return nextPageToken, nil
}

func (srv *ServiceImpl) StartFetchingCron(query string) {
	resultChan := make(chan []models.Video)
	var wg sync.WaitGroup
	var nextPageToken string
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			pageToken, err := fetchLatestVideos(query, nextPageToken, resultChan)
			if err != nil {
				log.Println("Error in calling Youtube API")
				return
			}
			nextPageToken = pageToken
			time.Sleep(10 * time.Second)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case videos := <-resultChan:
				{
					log.Println("Received", len(videos), "videos:", videos)
					log.Println("Dumping videos in elastic")
					err := srv.repo.InsertIntoElasticsearch(videos)
					if err != nil {
						log.Println("Error inserting into Elasticsearch:", err)
					}
				}
			default:
				continue
			}
		}
	}()
	wg.Wait()
}

func (srv *ServiceImpl) GetPaginatedResponse(pageSize, pageNumber int32) ([]models.Video, error) {
	return srv.repo.GetPaginatedResponse(pageSize, pageNumber)
}

func (srv *ServiceImpl) GetQueryResponse(query string) ([]models.Video, error) {
	return srv.repo.GetQueryResponse(query)
}
