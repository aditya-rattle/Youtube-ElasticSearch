package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"personal-project/fampay-assignment/models"
	"personal-project/fampay-assignment/utility"

	"github.com/elastic/go-elasticsearch/v8"
)

type RepositoryImpl struct {
	ESClient *elasticsearch.Client
}

func NewRepository(esClient *elasticsearch.Client) IRepository {
	return &RepositoryImpl{
		ESClient: esClient,
	}
}

func (repo *RepositoryImpl) InsertIntoElasticsearch(videos []models.Video) error {

	var buf bytes.Buffer

	for _, video := range videos {
		op := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "youtube_videos",
			},
		}

		if err := json.NewEncoder(&buf).Encode(op); err != nil {
			return fmt.Errorf("failed to encode index operation: %v", err)
		}

		if err := json.NewEncoder(&buf).Encode(video); err != nil {
			return fmt.Errorf("failed to encode video: %v", err)
		}
	}

	res, err := repo.ESClient.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return fmt.Errorf("failed to perform bulk insert: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.Status())
	}
	log.Printf("Bulk insert successful")
	return nil
}

func (repo *RepositoryImpl) GetPaginatedResponse(pageSize, pageNumber int32) ([]models.Video, error) {

	query := map[string]interface{}{
		"from": pageNumber * pageSize,
		"size": pageSize,
		"sort": []map[string]interface{}{
			{"publishedAt": "desc"},
		},
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		log.Println("[GetPaginatedResponse] error in marshalling queryjson", err)
		return nil, err
	}

	res, err := repo.ESClient.Search(
		repo.ESClient.Search.WithContext(context.Background()),
		repo.ESClient.Search.WithIndex("youtube_videos"),
		repo.ESClient.Search.WithBody(bytes.NewReader(queryJson)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.Status())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return utility.EntityMapper(response), nil
}

func (repo *RepositoryImpl) GetQueryResponse(queryString string) ([]models.Video, error) {

	requestBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": queryString,
				"fields": []string{
					"title",
					"description",
				},
			},
		},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request body: %v", err)
	}

	response, err := repo.ESClient.Search(
		repo.ESClient.Search.WithIndex("youtube_videos"),
		repo.ESClient.Search.WithBody(bytes.NewReader(requestBodyBytes)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search request: %v", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", response.Status())
	}

	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}
	return utility.EntityMapper(responseBody), nil
}
