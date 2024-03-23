package service

import "personal-project/fampay-assignment/models"

type IService interface {
	StartFetchingCron(query string)
	GetPaginatedResponse(pageSize, pageNumber int32) ([]models.Video, error)
	GetQueryResponse(query string) ([]models.Video, error)
}