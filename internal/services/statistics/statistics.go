package statistics

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

func GetActualStateByGeoJson() (dsApi.AreasResponse, error) {
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
	newRecords, err := fetchNewRecords(db)
	if err != nil {
		return 0, err
	}

	return dsRepo.SaveHistoryRecords(db, sliceutils.Map(newRecords, func(r dsApi.HistoryRecord, _ int, _ []dsApi.HistoryRecord) models.HistoryRecord {
		return recordToModel(r)
	})), nil
}

func GetMapHistory(db *gorm.DB) error {
	records, err := fetchNewRecords(db)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(len(records))

	cli := gentleman.New()

	geodata := make([]models.HistoryRecordGeoJSON, 0, len(records))
	for _, record := range records {
		go func(record dsApi.HistoryRecord) {
			defer wg.Done()
			areas, err := dsApi.GetHistoryRecordGeoJson(cli, record.Id)
			if err != nil {
				println(err.Error())
				return
			}
			mu.Lock()
			geodata = append(geodata, models.HistoryRecordGeoJSON{
				HistoryRecordId: record.Id,
				Geojson:         string(areas),
			})
			mu.Unlock()
		}(record)
	}
	wg.Wait()
	//TODO: Save geojson with history records, so as not to break fk constraints
	dsRepo.SaveGeoJson(db, geodata)

	return nil
}

func GetRecordsByDate(db *gorm.DB, t time.Time) []models.HistoryRecord {
	return dsRepo.FindRecordsByDate(db, t)
}

func fetchNewRecords(db *gorm.DB) ([]dsApi.HistoryRecord, error) {
	cli := gentleman.New()
	records, err := dsApi.GetHistoryRecords(cli)
	if err != nil {
		return nil, err
	}

	newRecordIds := dsRepo.FindNewIds[models.HistoryRecord](db, sliceutils.Pluck(records, func(r dsApi.HistoryRecord) *int { return &r.Id }))
	newRecords := sliceutils.Filter(records, func(r dsApi.HistoryRecord, _ int, _ []dsApi.HistoryRecord) bool {
		return sliceutils.Includes(newRecordIds, r.Id)
	})
	return newRecords, nil
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
