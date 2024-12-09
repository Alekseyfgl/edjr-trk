package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Paginate - common structure for pagination
type Paginate[T any] struct {
	PageNumber     int `json:"pageNumber"`
	RowTotalCount  int `json:"rowTotalCount"`
	TotalPageCount int `json:"totalPageCount"`
	PageSize       int `json:"pageSize"`
	Items          []T `json:"items"`
}

// RowArticle - структура для хранения данных статьи.
type RowArticle struct {
	ID    primitive.ObjectID `bson:"_id"`
	Text  string             `bson:"text"`
	Title string             `bson:"title"`
	Img   *string            `bson:"img"`
	Date  time.Time          `bson:"date"`
}

type ArticleResponse struct {
	ID    primitive.ObjectID `json:"id,omitempty"`
	Text  string             `json:"text,omitempty"`
	Title string             `json:"title,omitempty"`
	Img   *string            `json:"img"`
	Date  time.Time          `json:"date,omitempty"`
}

func (ar *RowArticle) CreateArtResp() *ArticleResponse {
	return &ArticleResponse{
		ID:    ar.ID,
		Title: ar.Title,
		Text:  ar.Text,
		Img:   ar.Img,
		Date:  ar.Date,
	}
}
