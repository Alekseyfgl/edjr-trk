package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RowProduct struct {
	ID    primitive.ObjectID `bson:"_id"`
	Text  string             `bson:"text"`
	Title string             `bson:"title"`
	Img   *string            `bson:"img"`
	Date  time.Time          `bson:"date"`
}

type ProductResponse struct {
	ID    primitive.ObjectID `json:"id,omitempty"`
	Text  string             `json:"text,omitempty"`
	Title string             `json:"title,omitempty"`
	Img   *string            `json:"img"`
	Date  time.Time          `json:"date,omitempty"`
}

func (ar *RowProduct) CreateProductResp() *ProductResponse {
	return &ProductResponse{
		ID:    ar.ID,
		Title: ar.Title,
		Text:  ar.Text,
		Img:   ar.Img,
		Date:  ar.Date,
	}
}
