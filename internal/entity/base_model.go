package entity

import (
	"time"

	"github.com/google/uuid"
)

type baseModel struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func initEntity() baseModel {
	now := time.Now()

	return baseModel{
		Id:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
