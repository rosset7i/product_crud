package entity

import (
	"time"

	"github.com/google/uuid"
)

type baseModel struct {
	Id         uuid.UUID
	created_at time.Time
	updated_at time.Time
}

func initEntity() baseModel {
	now := time.Now()

	return baseModel{
		Id:         uuid.New(),
		created_at: now,
		updated_at: now,
	}
}
