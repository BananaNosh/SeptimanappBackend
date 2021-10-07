package types

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	ID    int             `gorm:"primary_key, AUTO_INCREMENT"`
	Start time.Time       `validate:"required"`
	End   time.Time       `validate:"required,gtfield=Start"`
	Names []LocatedString `gorm:"polymorphic:Parent;" validate:"required"`
}

type Events []Event

func (event Event) MarshalJSON() ([]byte, error) {
	/**
	Marshal Event to json
	*/
	namesMap := make(map[string]string, len(event.Names))
	for _, name := range event.Names {
		namesMap[name.Language] = name.Value
	}
	return json.Marshal(struct {
		Id    int               `json:"id"`
		Start int64             `json:"start"`
		End   int64             `json:"end"`
		Names map[string]string `json:"names"`
	}{
		Id:    event.ID,
		Start: event.Start.Unix(),
		End:   event.End.Unix(),
		Names: namesMap,
	})
}

func (event *Event) UnmarshalJSON(data []byte) error {
	/**
	Marshal Event to json
	*/
	var aux struct {
		Id    int
		Start int64
		End   int64
		Names map[string]string
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	event.ID = aux.Id
	event.Start = time.Unix(aux.Start, 0)
	event.End = time.Unix(aux.End, 0)

	for language, name := range aux.Names {
		event.Names = append(event.Names, LocatedString{
			Value:    name,
			Language: language,
		})
	}
	return nil
}

type OverallLocation string

type Location struct {
	gorm.Model
	ID              string `gorm:"primary_key"`
	OverallLocation OverallLocation
	Longitude       float32
	Latitude        float32
	Altitude        float32
	IsMain          bool
	Titles          []LocatedString `gorm:"polymorphic:Parent;polymorphicValue:locations_title"`
	Descriptions    []LocatedString `gorm:"polymorphic:Parent;polymorphicValue:locations_description"`
}

type Language string

type LocatedString struct {
	gorm.Model
	ID         int `gorm:"primary_key, AUTO_INCREMENT" json:"-"`
	Value      string
	Language   string // TODO make lannguage (does not reall ymake a diffference
	ParentID   string `json:"-"`
	ParentType string `json:"-"`
}

func (location Location) MarshalJSON() ([]byte, error) {
	/**
	Marshal Event to json
	*/
	titlesMap := locationStringsToMap(location.Titles)
	descriptionsMap := locationStringsToMap(location.Descriptions)
	return json.Marshal(struct {
		ID              string            `json:"id"`
		OverallLocation OverallLocation   `json:"overallLocation"`
		Longitude       float32           `json:"longitude"`
		Latitude        float32           `json:"latitude"`
		Altitude        float32           `json:"altitude"`
		IsMain          bool              `json:"isMain"`
		Titles          map[string]string `json:"titles"`
		Descriptions    map[string]string `json:"descriptions"`
	}{
		ID:              location.ID,
		OverallLocation: location.OverallLocation,
		Longitude:       location.Longitude,
		Latitude:        location.Latitude,
		Altitude:        location.Altitude,
		IsMain:          location.IsMain,
		Titles:          titlesMap,
		Descriptions:    descriptionsMap,
	})
}

func (location *Location) UnmarshalJSON(data []byte) (err error) { // TODO add method matching Marshal -> typedef
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

func locationStringsToMap(strings []LocatedString) map[string]string {
	stringMap := make(map[string]string, len(strings))
	for _, title := range strings {
		stringMap[title.Language] = title.Value
	}
	return stringMap
}

type ApiKeyInfo struct {
	ApiKeyHash string `gorm:"primary_key,column:api_key"`
}
