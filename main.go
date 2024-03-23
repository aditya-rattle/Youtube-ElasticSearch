package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"personal-project/fampay-assignment/client"
	"personal-project/fampay-assignment/controller"
	"personal-project/fampay-assignment/repo"
	"personal-project/fampay-assignment/service"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var ESClient *elasticsearch.Client

func createIndicesIfNotExist() error {
	// Define the index mappings
	indexName := "youtube_videos"
	mappings := `{
		"settings": {
			"index.max_ngram_diff": 5,
			"analysis": {
				"tokenizer": {
					"name_ngram_tokenizer": {
						"type":        "ngram",
						"min_gram":    3,
						"max_gram":    6,
						"token_chars": ["letter", "digit"]
					}
				},
				"analyzer": {
					"name_ngram_analyzer": {
						"type":      "custom",
						"tokenizer": "name_ngram_tokenizer",
						"filter":    ["lowercase", "asciifolding"]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"title": {
					"type": "text"
				},
				"description": {
					"type": "text"
				},
				"publishedAt": {
					"type": "date"
				},
				"thumbnails": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				}
			}
		}
	}`

	res, err := ESClient.Indices.Exists([]string{indexName})
	if err != nil {
		log.Println("Error in checking indices", err)
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		res, err := ESClient.Indices.Create(indexName,
			ESClient.Indices.Create.WithBody(strings.NewReader(mappings)))
		if err != nil {
			log.Println("Error in creating indices", err)
			return err
		}
		defer res.Body.Close()

		// Check the response status
		if res.IsError() {
			log.Printf(fmt.Sprintf("error creating index: %s", res.Status()))
			return fmt.Errorf("error in creating indices %s", res.Status())
		}

		log.Println("Index created successfully")

		return nil
	}
	log.Println("Indices already exists")
	return nil
}

func initElasticSearch() error {
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ESHost")},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   1000,
			ResponseHeaderTimeout: time.Second * time.Duration(3),
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 25 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Username: os.Getenv("ESUsername"),
		Password: os.Getenv("ESPassword"),
		APIKey:   os.Getenv("ESAPIKey"),
	}

	var err error
	ESClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error in connecting to Elasticsearch: %v", err)
	}

	res, err := ESClient.Ping()
	if err != nil {
		log.Fatalf("Error pinging Elasticsearch server: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Elasticsearch server returned an error: %s", res.String())
	}

	log.Println("Successfully connected to Elasticsearch")

	if err := createIndicesIfNotExist(); err != nil {
		return fmt.Errorf("error in creating Elasticsearch indices: %v", err)
	}

	log.Println("Successfully created Elasticsearch indices")

	return nil
}

func main() {
	fmt.Println("Hello, Fam!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = initElasticSearch()
	if err != nil {
		panic("Error in starting server")
	}
	apiKeysStr := os.Getenv("YTAPIKEY")
	if apiKeysStr == "" {
		log.Fatal("YTAPIKEYS environment variable not set")
	}
	// Split the string into a slice of strings
	apiKeys := strings.Split(apiKeysStr, ",")

	youTubeClient := client.NewYouTubeAPIClient(apiKeys)
	repo := repo.NewRepository(ESClient)
	service := service.NewService(repo, youTubeClient)
	controller := controller.NewController(service)
	router := mux.NewRouter()
	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/start-fetching-videos", controller.FetchVideos).Methods(http.MethodPost)
	apiV1.HandleFunc("/get-videos", controller.GetVideos).Methods(http.MethodGet)
	apiV1.HandleFunc("/get-query-videos", controller.GetVideosByQueryParams).Methods(http.MethodGet)
	log.Fatal("Server starting at port: 9000", http.ListenAndServe(":9000", router))

}
