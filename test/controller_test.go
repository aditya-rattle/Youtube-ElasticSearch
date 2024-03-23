package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-project/fampay-assignment/controller"
	"personal-project/fampay-assignment/models"
	"personal-project/fampay-assignment/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestController_FetchVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)
	controller := controller.NewController(mockService)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/start-fetching-videos?query=test", nil)
	w := httptest.NewRecorder()
	mockService.EXPECT().StartFetchingCron("test")
	controller.FetchVideos(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestController_GetVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)
	controller := controller.NewController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/get-videos?pageNumber=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	mockVideos := []models.Video{{Id: "1", Title: "Video 1"}, {Id: "2", Title: "Video 2"}}
	mockService.EXPECT().GetPaginatedResponse(int32(10), int32(1)).Return(mockVideos, nil)

	controller.GetVideos(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody, _ := json.Marshal(mockVideos)
	assert.Equal(t, expectedBody, w.Body.Bytes())
}

func TestController_GetVideosByQueryParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)
	controller := controller.NewController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/get-query-videos?queryString=test", nil)
	w := httptest.NewRecorder()
	mockVideos := []models.Video{{Id: "1", Title: "Video 1"}, {Id: "2", Title: "Video 2"}}
	mockService.EXPECT().GetQueryResponse("test").Return(mockVideos, nil)
	controller.GetVideosByQueryParams(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody, _ := json.Marshal(mockVideos)
	assert.Equal(t, expectedBody, w.Body.Bytes())
}
