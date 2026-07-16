package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"reflect"
	"testing"
	"time"
)

// the latin horarium of data/testData/horariumHelper
func testHorarium() Horarium {
	return Horarium{
		Events: []WeekViewEvent{
			{"e0", eventTime1, eventTime2, "proba0"},
			{"ev1", eventTime3, eventTime4, "proba1"},
		},
		Language: "la",
	}
}

func TestEventTime_ToTime(t *testing.T) {
	// The wanted times are spelled out instead of being built with ToTime:
	// an expectation that goes through ToTime itself cannot catch a wrong ToTime.
	berlin := util.Locale()
	tests := []struct {
		name     string
		evTime   EventTime
		location *time.Location
		want     time.Time
	}{
		{"start of test event", EventTime{Year: 2019, Month: 6, Day: 25, Hour: 16, Minute: 30},
			berlin, time.Date(2019, time.July, 25, 16, 30, 0, 0, berlin)},
		{"full hour keeps minute 0", EventTime{Year: 2019, Month: 6, Day: 25, Hour: 18, Minute: 0},
			berlin, time.Date(2019, time.July, 25, 18, 0, 0, 0, berlin)},
		{"quarter past", EventTime{Year: 2020, Month: 7, Day: 3, Hour: 9, Minute: 15},
			berlin, time.Date(2020, time.August, 3, 9, 15, 0, 0, berlin)},
		{"january is month 0", EventTime{Year: 2020, Month: 0, Day: 1, Hour: 0, Minute: 0},
			berlin, time.Date(2020, time.January, 1, 0, 0, 0, 0, berlin)},
		{"december is month 11", EventTime{Year: 2020, Month: 11, Day: 31, Hour: 23, Minute: 59},
			berlin, time.Date(2020, time.December, 31, 23, 59, 0, 0, berlin)},
		{"minute differs from month", EventTime{Year: 2021, Month: 3, Day: 15, Hour: 9, Minute: 45},
			berlin, time.Date(2021, time.April, 15, 9, 45, 0, 0, berlin)},
		{"other location", EventTime{Year: 2019, Month: 6, Day: 25, Hour: 16, Minute: 30},
			time.UTC, time.Date(2019, time.July, 25, 16, 30, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.evTime.ToTime(tt.location)
			if !got.Equal(tt.want) || got.Location() != tt.want.Location() {
				t.Errorf("ToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHorarium_toEventList(t *testing.T) {
	// a horarium holds one language, so only the latin name of each event survives
	latinOnly := func() []types.Event {
		wanted := testEvents()
		for i := range wanted {
			wanted[i].Names = wanted[i].Names[1:2]
		}
		return wanted
	}
	tests := []struct {
		name     string
		horarium Horarium
		idOffset int
		want     []types.Event
	}{
		{"empty horarium", Horarium{Events: []WeekViewEvent{}, Language: "la"}, 0, []types.Event{}},
		{"horarium1", testHorarium(), 0, latinOnly()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			horarium := Horarium{
				Events:   tt.horarium.Events,
				Language: tt.horarium.Language,
			}
			events := horarium.ToEventList(tt.idOffset)
			for i := range events {
				events[i].Model.CreatedAt = time.Time{}
				events[i].Model.UpdatedAt = time.Time{}
			}
			if got := events; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toEventList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventsFromJsonHoraria(t *testing.T) {
	// the events as they come out of the json, shifted by the id offset
	wantedFromOffset := func(offset int) []types.Event {
		wanted := testEvents()
		for i := range wanted {
			wanted[i].ID = i + offset
			wanted[i].Model.ID = uint(i + offset)
		}
		return wanted
	}
	tests := []struct {
		name         string
		dataPath     string
		dataIdOffset int
		want         []types.Event
	}{
		{"test1", "../data/testData/horariumHelper/", 10, wantedFromOffset(10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := EventsFromJsonHoraria(tt.dataPath, tt.dataIdOffset)
			for i := range events {
				events[i].Model.CreatedAt = time.Time{}
				events[i].Model.UpdatedAt = time.Time{}
			}
			if got := events; !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("eventsFromJsonHoraria() = %v, want %v", got, tt.want)
			}
		})
	}
}
