package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Album struct {
    Id primitive.ObjectID `json:"id"`
    Title string `json:"title" validate:"required"`
    Artist string `json:"artist" validate:"required"`
    Price float64 `json:"price" validate:"required"`
}