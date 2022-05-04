package models

import (
	"strconv"
	"time"
)

type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;not null"`

	CreatedAt time.Time `gorm:"column:createdAt;index"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index;after createdat"`
}

func (a BaseModel) GetStringId() string {
	return strconv.FormatUint(a.ID, 10)
}
