package models

import "time"

type HistoryRecord struct {
	Id            int `gorm:"primaryKey;autoIncrement:false"`
	DescriptionEn string
	Description   string
	CreatedAtDS   time.Time
	UpdatedAtDS   time.Time
	CreatedAt     time.Time
	Datetime      string
	Status        bool
}
