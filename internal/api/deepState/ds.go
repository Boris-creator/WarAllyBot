package deepstate

import (
	"ds/internal/api"
	"encoding/json"
	"fmt"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/url"
)

type HistoryRecord struct {
	Id            int    `json:"id"`
	DescriptionEn string `json:"descriptionEn"`
	Description   string `json:"description"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	Datetime      string `json:"datetime"`
	Status        bool   `json:"status"`
}

type LastRecordResponse struct {
	Id       int    `json:"id"`
	Datetime string `json:"datetime"`
	Map      any    `json:"map"`
}

type AreaStatusType string

const (
	Unspecified       AreaStatusType = "unspecified"
	Liberated         AreaStatusType = "liberated"
	Occupied_to       AreaStatusType = "occupied_to_24_02_2022"
	Occupied_after    AreaStatusType = "occupied_after_24_02_2022"
	Occupied          AreaStatusType = "occupied"
	Other_territories AreaStatusType = "other_territories"
)

type Area struct {
	Type    AreaStatusType `json:"type"`
	Area    float32        `json:"area"`
	Percent string         `json:"percent"`
	Hash    string         `json:"hash"`
}
type AreasResponse []Area

type GeojsonResponse interface{}

const baseUrl = "https://deepstatemap.live"

func GetLastHistoryRecord(cli *gentleman.Client) (*LastRecordResponse, error) {
	cli.BaseURL(baseUrl)
	response, err := cli.Request().Path("/api/history/last").Send()
	err = api.HandleError(*response, err)
	if err != nil {
		return nil, err
	}
	var r LastRecordResponse
	json.Unmarshal(response.Bytes(), &r)
	return &r, nil
}

func GetHistoryRecords(cli *gentleman.Client) ([]HistoryRecord, error) {
	cli.BaseURL(baseUrl).Use(url.Path("/api/history"))
	response, err := cli.Get().Send()
	err = api.HandleError(*response, err)
	if err != nil {
		return nil, err
	}
	var rs []HistoryRecord
	json.Unmarshal(response.Bytes(), &rs)
	return rs, nil
}

func GetHistoryRecordAreas(cli *gentleman.Client, recordId int) (*AreasResponse, error) {
	cli.BaseURL(baseUrl).Use(url.Path("/api/history/:id/areas"))
	cli.Use(url.Param("id", fmt.Sprintf("%d", recordId)))
	response, err := cli.Get().Send()
	err = api.HandleError(*response, err)
	if err != nil {
		return nil, err
	}
	var r AreasResponse
	json.Unmarshal(response.Bytes(), &r)
	return &r, nil
}

func GetHistoryRecordGeoJson(cli *gentleman.Client, recordId int) ([]byte, error) {
	cli.BaseURL(baseUrl).Use(url.Path("/api/history/:id/geojson"))
	cli.Use(url.Param("id", fmt.Sprintf("%d", recordId)))
	response, err := cli.Get().Send()
	err = api.HandleError(*response, err)
	if err != nil {
		return nil, err
	}
	return response.Bytes(), nil
}
