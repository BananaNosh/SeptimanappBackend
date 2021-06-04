package types

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	ID    int `gorm:"primary_key, AUTO_INCREMENT"`
	Start time.Time
	End   time.Time
	Names []LocatedString `gorm:"polymorphic:Parent;"`
}

func (event Event) MarshalJSON() ([]byte, error) {
	/**
	Unmarshal json bytes to location
	*/
	namesMap := make(map[string]string, len(event.Names))
	for _, name := range event.Names {
		namesMap[name.Language] = name.Value
	}
	return json.Marshal(struct {
		Id    int
		Start int64
		End   int64
		Names map[string]string
	}{
		Id:    event.ID,
		Start: event.Start.Unix(),
		End:   event.End.Unix(),
		Names: namesMap,
	})
}

type Location struct {
	gorm.Model
	ID              string `gorm:"primary_key"`
	OverallLocation string
	Longitude       float32
	Latitude        float32
	Altitude        float32
	IsMain          bool
	Titles          []LocatedString `gorm:"polymorphic:Parent;"`
	Descriptions    []LocatedString `gorm:"polymorphic:Parent;"`
}

type LocatedString struct {
	gorm.Model
	ID         int `gorm:"primary_key, AUTO_INCREMENT" json:"-"`
	Value      string
	Language   string
	ParentID   string `json:"-"`
	ParentType string `json:"-"`
}

type Language string

func (location *Location) UnmarshalJSON(data []byte) (err error) {
	/**
	Unmarshal json bytes to location
	*/
	var auxiliaryLocation struct {
		Id             string
		TitleMap       map[string]string
		DescriptionMap map[string]string
		IsMain         bool
		Coordinates    map[string]float32
	}
	if err = json.Unmarshal(data, &auxiliaryLocation); err == nil {
		location.ID = auxiliaryLocation.Id
		location.Longitude = auxiliaryLocation.Coordinates["mLongitude"]
		location.Latitude = auxiliaryLocation.Coordinates["mLatitude"]
		location.Altitude = auxiliaryLocation.Coordinates["mAltitude"]
		location.IsMain = auxiliaryLocation.IsMain
		location.Titles = locationStringsFromMap(auxiliaryLocation.TitleMap)
		location.Descriptions = locationStringsFromMap(auxiliaryLocation.DescriptionMap)
	}

	return err
}

func locationStringsFromMap(stringMap map[string]string) []LocatedString {
	var locationStrings []LocatedString
	for k, v := range stringMap {
		locationStrings = append(locationStrings, LocatedString{
			Model:    gorm.Model{},
			Value:    v,
			Language: k,
		})
	}
	return locationStrings
}
