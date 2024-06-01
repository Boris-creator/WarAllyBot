package models

import (
	"time"

	"gorm.io/gorm"
)

type HistoryRecordGeoJSON struct {
	gorm.Model
	HistoryRecordId int `gorm:"unique"`
	HistoryRecord   HistoryRecord
	Geojson         string `gorm:"type:longtext"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
