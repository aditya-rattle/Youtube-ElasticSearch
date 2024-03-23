package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"personal-project/fampay-assignment/models"
	"personal-project/fampay-assignment/service"
	"strconv"
)

type Controller struct {
	service service.IService
}

func NewController(service service.IService) *Controller {
	return &Controller{
		service: service,
	}
}

// Fetch videos from YouTube and store in ES.
func (ctr *Controller) FetchVideos(w http.ResponseWriter, r *http.Request) {

	log.Println("Fetching Videos")
	queryParams := r.URL.Query()
	queryString := queryParams.Get("query")
	if len(queryString) == 0 {
		log.Println("No query string provided")
		w.WriteHeader(http.StatusBadGateway)
	}
	ctr.service.StartFetchingCron(queryString)

}

// Return paginated response in descending order.
func (ctr *Controller) GetVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	queryParams := r.URL.Query()
	pageSize, err := strconv.ParseInt(queryParams.Get("pageSize"), 10, 32)
	if err != nil {
		log.Println("error in parsing int pageSize", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	pageNum, err := strconv.ParseInt(queryParams.Get("pageNumber"), 10, 32)
	if err != nil {
		log.Println("error in parsing int pageNum", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if pageNum == 0 || pageSize == 0 {
		log.Println("invalid pagenum and pagesize ")
		w.WriteHeader(http.StatusBadRequest)
	}
	videos, err = ctr.service.GetPaginatedResponse(int32(pageSize), int32(pageNum))
	if err != nil {
		log.Println("Error in GetPaginatedResponse", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	videosB, err := json.Marshal(videos)
	if err != nil {
		log.Println("Error in Marshaling paginated", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(videosB)
}

// Return all matching response with given query(description, title).
func (ctr *Controller) GetVideosByQueryParams(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	var err error
	queryParams := r.URL.Query()
	queryString := queryParams.Get("queryString")
	if len(queryString) == 0 {
		log.Println("invalid query string")
		w.WriteHeader(http.StatusBadRequest)
	}
	videos, err = ctr.service.GetQueryResponse(queryString)
	if err != nil {
		log.Println("error in GetQueryResponse", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	videosB, err := json.Marshal(videos)
	if err != nil {
		log.Println("Error in Marshaling paginated", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(videosB)	
}
