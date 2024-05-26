package deepstate

import (
	dsApi "ds/internal/api/deepState"
	"ds/internal/database/models"
	dsRepo "ds/internal/database/repositories"
	"sync"
	"time"

	"github.com/Goldziher/go-utils/sliceutils"
	"gopkg.in/h2non/gentleman.v2"
	"gorm.io/gorm"
)

func GetActualState() (dsApi.AreasResponse, error) {
	cli := gentleman.New()
	rec, err := dsApi.GetLastHistoryRecord(cli)
	if err != nil {
		return dsApi.AreasResponse{}, err
	}
	areas, err := dsApi.GetHistoryRecordAreas(cli, rec.Id)
	if err != nil {
		return dsApi.AreasResponse{}, err
	}
	return *areas, nil
}

func GetHistory(db *gorm.DB) (int, error) {
	cli := gentleman.New()
	records, err := dsApi.GetHistoryRecords(cli)
	if err != nil {
		return 0, err
	}
	newRecordIds := dsRepo.FindNewIds(db, sliceutils.Pluck(records, func(r dsApi.HistoryRecord) *int { return &r.Id }))
	newRecords := sliceutils.Filter(records, func(r dsApi.HistoryRecord, _ int, _ []dsApi.HistoryRecord) bool {
		return sliceutils.Includes(newRecordIds, r.Id)
	})
	return dsRepo.SaveHistoryRecords(db, sliceutils.Map(newRecords, func(r dsApi.HistoryRecord, _ int, _ []dsApi.HistoryRecord) models.HistoryRecord {
		return recordToModel(r)
	})), nil
}

func GetAreasControlHistory() error {
	cli := gentleman.New()
	records, err := dsApi.GetHistoryRecords(cli)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(len(records))

	areasRecords := make([]dsApi.AreasResponse, len(records))
	for _, record := range records {
		go func(record dsApi.HistoryRecord) {
			areas, err := dsApi.GetHistoryRecordAreas(cli, record.Id)
			if err != nil {
				return
			}
			defer func() {
				wg.Done()
			}()
			mu.Lock()
			areasRecords = append(areasRecords, *areas)
			mu.Unlock()
		}(record)
	}
	wg.Wait()

	return nil
}

func GetRecordsByDate(db *gorm.DB, t time.Time) []models.HistoryRecord {
	return dsRepo.FindRecordsByDate(db, t)
}

func recordToModel(r dsApi.HistoryRecord) models.HistoryRecord {
	tc, _ := time.Parse("2006-01-02T15:04:05Z", r.CreatedAt)
	uc, _ := time.Parse("2006-01-02T15:04:05Z", r.UpdatedAt)
	return models.HistoryRecord{
		Id:            r.Id,
		Description:   r.Description,
		DescriptionEn: r.DescriptionEn,
		CreatedAtDS:   tc,
		UpdatedAtDS:   uc,
		Status:        r.Status,
		Datetime:      r.Datetime,
	}
}
