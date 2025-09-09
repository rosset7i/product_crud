package domain

import (
	"time"

	"github.com/google/uuid"
)

type baseModel struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func initEntity() baseModel {
	now := time.Now()

	return baseModel{
		Id:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
