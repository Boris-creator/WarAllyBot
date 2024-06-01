package repositories

import (
	"ds/internal/database/models"
	"time"

	"github.com/Goldziher/go-utils/sliceutils"
	"gorm.io/gorm"
)

type Model interface {
	models.HistoryRecord | models.HistoryRecordGeoJSON
}

func FindNewIds[M Model](db *gorm.DB, ids []int) []int {
	var existing []int
	var model M
	db.Model(&model).Where(ids).Pluck("id", &existing)
	return sliceutils.Filter(ids, func(id int, _ int, _ []int) bool {
		return !sliceutils.Includes(existing, id)
	})
}

func FindRecordsByDate(db *gorm.DB, t time.Time) []models.HistoryRecord {
	var res []models.HistoryRecord

	formatStr := "2006-01-02 00:00:00"
	dateFrom := t.Format(formatStr)
	dateTo := t.AddDate(0, 0, 1).Format(formatStr)

	db.Where("created_at_ds > ?", dateFrom).Where("created_at_ds < ?", dateTo).Find(&res)
	return res
}

func GetlastRecord(db *gorm.DB) models.HistoryRecord {
	var record models.HistoryRecord
	db.Order("created_at_ds desc").First(&record)
	return record
}

func SaveHistoryRecords(db *gorm.DB, records []models.HistoryRecord) int {
	tx := db.CreateInBatches(records, len(records))
	return int(tx.RowsAffected)
}

func SaveGeoJson(db *gorm.DB, records []models.HistoryRecordGeoJSON) int {
	tx := db.CreateInBatches(records, len(records))
	return int(tx.RowsAffected)
}
