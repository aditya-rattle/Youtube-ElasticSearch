package repo

import "personal-project/fampay-assignment/models"

type IRepository interface {
	InsertIntoElasticsearch([]models.Video) error
	GetPaginatedResponse(pageSize, pageNumber int32) ([]models.Video, error)
	GetQueryResponse(query string) ([]models.Video, error)
}