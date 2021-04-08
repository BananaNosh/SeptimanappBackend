package main

import (
	"encoding/json"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"
)

const horariumFilePattern = "horarium_\\d+_(?P<lang>\\w+).json"

type Horarium struct {
	Events   []WeekViewEvent `json:"events"`
	Language string
}

func (horarium Horarium) ToEventList(idOffset int) []Event {
	var events = make([]Event, 0)
	for i, event := range horarium.Events {
		id := idOffset + i
		events = append(events, Event{
			Model:    gorm.Model{ID: uint(id), CreatedAt: time.Now(), UpdatedAt: time.Now()},
			ID:       id,
			Start:    event.StartTime.ToTime(locale),
			End:      event.EndTime.ToTime(locale),
			Name:     event.Name,
			Language: horarium.Language,
		})
	}
	return events
}

type WeekViewEvent struct {
	ID        string    `json:"mId"`
	StartTime EventTime `json:"mStartTime"`
	EndTime   EventTime `json:"mEndTime"`
	Name      string    `json:"mName"`
}

type EventTime struct {
	Year   int `json:"year"`
	Month  int `json:"month"` // 0-11
	Day    int `json:"dayOfMonth"`
	Hour   int `json:"hourOfDay"`
	Minute int `json:"minute"`
}

func (evTime EventTime) ToTime(location *time.Location) time.Time {
	return time.Date(evTime.Year, time.Month(evTime.Month+1), evTime.Day, evTime.Hour, evTime.Month, 0, 0, location)
}

func EventsFromJsonHoraria(dataPath string, dataIdOffset int) []Event {
	var allEvents []Event

	// Call Readdir to get all files.
	outputDirFiles, _ := ioutil.ReadDir(dataPath)

	// compile regex for HorariaFiles
	reg := regexp.MustCompile(horariumFilePattern)
	offset := dataIdOffset
	for _, file := range outputDirFiles {
		match := reg.FindStringSubmatch(file.Name())
		isHorariumFile := len(match) > 1
		if isHorariumFile {
			language := match[1]

			// open the file pointer
			filePath := path.Join(dataPath, file.Name())

			if horarium, err := readHorariumFromFile(filePath); err == nil {
				events := horarium.ToEventList(offset)
				offset += len(events)
				// set correct language
				for i := range events {
					events[i].Language = language
				}
				allEvents = append(allEvents, events...)
			}
		}
	}
	return allEvents
}

func readHorariumFromFile(filePath string) (Horarium, error) {
	var horarium Horarium
	if horariumFile, err := os.Open(filePath); err == nil {
		defer horariumFile.Close()

		// create a new decoder
		err = json.NewDecoder(horariumFile).Decode(&horarium)
	}
	return horarium, nil
}
