package database

import (
	"SeptimanappBackend/types"
	"reflect"
	"sort"
	"testing"
)

var locations []types.Location

func init() {
	location1 := types.Location{
		ID:              "test1",
		OverallLocation: "over1",
		Longitude:       100,
		Latitude:        1000,
		Altitude:        10,
		IsMain:          true,
		Titles: []types.LocatedString{
			{Value: "testTitle1de", Language: "de"},
			{Value: "testTitle1la", Language: "la"},
		},
		Descriptions: []types.LocatedString{
			{Value: "testDesc1de", Language: "de"},
			{Value: "testDesc1la", Language: "la"},
		},
	}
	location2 := types.Location{
		ID:              "test2",
		OverallLocation: "over1",
		Longitude:       101,
		Latitude:        1001,
		Altitude:        1,
		IsMain:          true,
		Titles: []types.LocatedString{
			{Value: "testTitle2de", Language: "de"},
			{Value: "testTitle2la", Language: "la"},
		},
		Descriptions: []types.LocatedString{
			{Value: "testDesc2de", Language: "de"},
			{Value: "testDesc2la", Language: "la"},
		},
	}
	locations = append(locations, location1, location2)
}

func TestLocationsFromJsonFiles(t *testing.T) {
	tests := []struct {
		name     string
		dataPath string
		want     []types.Location
	}{
		{
			"test1",
			"../data/testData/locationHelper/",
			locations,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LocationsFromJsonFiles(tt.dataPath)
			for _, location := range got {
				sort.Slice(location.Titles, func(i, j int) bool {
					return location.Titles[i].Language < location.Titles[j].Language
				})
				sort.Slice(location.Descriptions, func(i, j int) bool {
					return location.Descriptions[i].Language < location.Descriptions[j].Language
				})
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocationsFromJsonFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
