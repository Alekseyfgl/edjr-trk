package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RowProduct struct {
	ID        primitive.ObjectID `bson:"_id"`
	Text      string             `bson:"text"`
	ShortText string             `bson:"shortText"`
	Title     string             `bson:"title"`
	Img       *string            `bson:"img"`
	Date      time.Time          `bson:"date"`
}

type ProductResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Text      string             `json:"text,omitempty"`
	ShortText string             `json:"shortText"`
	Title     string             `json:"title,omitempty"`
	Img       *string            `json:"img"`
	Date      time.Time          `json:"date,omitempty"`
}

func (ar *RowProduct) CreateProductResp() *ProductResponse {
	return &ProductResponse{
		ID:        ar.ID,
		Title:     ar.Title,
		ShortText: ar.ShortText,
		Text:      ar.Text,
		Img:       ar.Img,
		Date:      ar.Date,
	}
}
