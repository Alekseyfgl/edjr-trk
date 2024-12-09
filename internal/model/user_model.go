package model

import (
	"edjr-trk/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// RowUser - структура для хранения данных user.
type RowUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Phone     string             `bson:"phone"`
	IsAdmin   bool               `bson:"IsAdmin"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

// UserResponse - for UI response
type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	IsAdmin   bool               `json:"isAdmin"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func (u *RowUser) CreateUserResp() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *RowUser) HashPassword() error {
	hashedPassword, err := utils.HashData(u.Password, 12)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}
