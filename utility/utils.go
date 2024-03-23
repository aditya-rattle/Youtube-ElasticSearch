package utility

import (
	"personal-project/fampay-assignment/models"
	"time"
)

func Mapper(entity []interface{}) []string {
	res := make([]string, 0)
	for _, v := range entity {
		res = append(res, v.(string))
	}
	return res
}


func EntityMapper(entity map[string]interface{}) []models.Video {
	videos := make([]models.Video, 0)
	hits := entity["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		title := source["title"].(string)
		description := source["description"].(string)
		publishedAt, _ := time.Parse(time.RFC3339, source["publishedAt"].(string))
		videoId := source["Id"].(string)
		videoURL := source["url"].(string)
		thumbnails := source["thumbnails"].([]interface{})

		videos = append(videos, models.Video{
			Title:       title,
			Description: description,
			PublishedAt: publishedAt,
			Thumbnails:  Mapper(thumbnails),
			Id:          videoId,
			URL:         videoURL,
		})
	}
	return videos
}