package service

import (
	"fmt"
	"log"
	"personal-project/fampay-assignment/client"
	"personal-project/fampay-assignment/models"
	"personal-project/fampay-assignment/repo"
	"sync"
	"time"
)

type ServiceImpl struct {
	repo   repo.IRepository
	client *client.YouTubeAPIClient
}

func NewService(repo repo.IRepository, client *client.YouTubeAPIClient) IService {
	return &ServiceImpl{
		repo:   repo,
		client: client,
	}
}

func (srv *ServiceImpl) fetchLatestVideos(query, pageToken string, resultChan chan<- []models.Video) (string, error) {

	response, err := srv.client.GetYoutubeVideoDetails(query, pageToken)
	if err != nil {
		log.Println("Error from youtube client", err)
		return "", err
	}
	var nextPageToken string
	videos := make([]models.Video, 0)
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
	/*
		Adding a routine which will fetch data
		from YT in async and put it in resultChan
	*/
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			pageToken, err := srv.fetchLatestVideos(query, nextPageToken, resultChan)
			if err != nil {
				log.Println("Error in calling Youtube API")
				return
			}
			nextPageToken = pageToken
			time.Sleep(10 * time.Second)
		}
	}()

	/*
		Adding a routine which will fetch data from
		resultChan in async and put it in ES
	*/
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
