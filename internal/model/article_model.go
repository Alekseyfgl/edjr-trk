package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// RowArticle - структура для хранения данных статьи.
type RowArticle struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Text  string             `json:"text,omitempty" bson:"text,omitempty"`
	Title string             `json:"title,omitempty" bson:"title,omitempty"`
	Img   *string            `json:"img,omitempty" bson:"img,omitempty"`
	Date  time.Time          `json:"date,omitempty" bson:"date,omitempty"`
}

type ArticleResponse struct {
	ID    primitive.ObjectID `json:"id,omitempty"`
	Text  string             `json:"text,omitempty"`
	Title string             `json:"title,omitempty"`
	Img   *string            `json:"img"`
	Date  time.Time          `json:"date,omitempty"`
}

func (ar *RowArticle) CreateArtResp() ArticleResponse {
	return ArticleResponse{
		ID:    ar.ID,
		Title: ar.Title,
		Text:  ar.Text,
		Img:   ar.Img,
		Date:  ar.Date,
	}
}
