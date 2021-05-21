package types

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	ID       int `gorm:"primary_key, AUTO_INCREMENT"`
	Start    time.Time
	End      time.Time
	Name     string
	Language string `json:"-"`
}

type Location struct {
	gorm.Model
	ID              string `gorm:"primary_key"`
	OverallLocation string
	Longitude       float32
	Latitude        float32
	Altitude        float32
	IsMain          bool
	Titles          []LocationString `gorm:"foreignKey:LocationID"`
	Descriptions    []LocationString `gorm:"foreignKey:LocationID"`
}

type LocationString struct {
	gorm.Model
	ID         int `gorm:"primary_key, AUTO_INCREMENT" json:"-"`
	Value      string
	Language   string
	LocationID string `json:"-"`
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

func locationStringsFromMap(stringMap map[string]string) []LocationString {
	var locationStrings []LocationString
	for k, v := range stringMap {
		locationStrings = append(locationStrings, LocationString{
			Model:    gorm.Model{},
			Value:    v,
			Language: k,
		})
	}
	return locationStrings
}
