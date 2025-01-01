package geography

import (
	"strings"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
)

type AreaControl struct {
	Liberated      float64
	OccupiedBefore float64
	OccupiedAfter  float64
	Unspecified    float64
}

// DeepState api provides a special endpoint to get consise data about areas control:
// https://deepstatemap.live/api/history/:id/areas
// but we can calculate it from geojson the next way
func CalculateAreas(d []byte) (AreaControl, error) {
	var ac AreaControl

	fc, err := geojson.UnmarshalFeatureCollection(d)
	if err != nil {
		return ac, err
	}

	for _, f := range fc.Features {
		if f.Geometry.GeoJSONType() != "Polygon" {
			continue
		}
		name, ok := f.Properties["name"]
		if !ok {
			continue
		}

		areaName := name.(string)
		delimeter := "///"
		status := areaName[strings.LastIndex(areaName, delimeter)+len(delimeter):]

		a := geo.Area(f.Geometry)
		status = strings.TrimSpace(status)
		if strings.HasPrefix(status, "geoJSON.status.dismissed") {
			ac.Liberated += a
		}
		switch status {
		case "geoJSON.status.occupied":
			ac.OccupiedAfter += a
		case "geoJSON.territories.ordlo":
			ac.OccupiedBefore += a
		case "geoJSON.territories.crimea":
			ac.OccupiedBefore += a
		case "geoJSON.status.unknown":
			ac.Unspecified += a
		}
	}
	return ac, nil
}
