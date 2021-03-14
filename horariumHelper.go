package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Horarium struct {
	Events   []WeekViewEvent `json:"events"`
	Language string
}

func (horarium Horarium) toEventList(idOffset int) []Event {
	var events []Event
	for i, event := range horarium.Events {
		id := idOffset + i
		events = append(events, Event{
			Model:    gorm.Model{ID: uint(i), CreatedAt: time.Now(), UpdatedAt: time.Now()},
			ID:       id,
			Start:    event.StartTime.toTime(locale),
			End:      event.EndTime.toTime(locale),
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

func (evTime EventTime) toTime(location *time.Location) time.Time {
	return time.Date(evTime.Year, time.Month(evTime.Month+1), evTime.Day, evTime.Hour, evTime.Month, 0, 0, location)
}

func insertEventsFromJsonHorarium(path string, db *gorm.DB) {
	// open the file pointer
	if horariumFile, err := os.Open(path); err == nil {
		defer horariumFile.Close()

		// initialize the storage for the decoded data
		var horarium Horarium

		// create a new decoder
		err := json.NewDecoder(horariumFile).Decode(&horarium)
		if err != nil {
			log.Fatal(err)
		}

		events := horarium.toEventList(0)
		for i := 0; i < 10; i++ {
			fmt.Println(events[i])
		}
	}
}
