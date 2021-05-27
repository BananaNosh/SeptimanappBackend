package database

import (
	"SeptimanappBackend/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

var eventTime1 EventTime
var eventTime2 EventTime
var eventTime3 EventTime
var eventTime4 EventTime
var horarium Horarium
var events []types.Event

func init() {
	eventTime1 = EventTime{
		Year:   2019,
		Month:  6,
		Day:    25,
		Hour:   16,
		Minute: 30,
	}
	eventTime2 = EventTime{
		Year:   2019,
		Month:  6,
		Day:    25,
		Hour:   18,
		Minute: 00,
	}
	eventTime3 = eventTime1
	eventTime3.Hour = 20
	eventTime4 = eventTime2
	eventTime4.Hour = 21
	var locale, _ = time.LoadLocation("Europe/Berlin")
	horariumEvents := []WeekViewEvent{
		{"e0", eventTime1, eventTime2, "proba0"},
		{"ev1", eventTime3, eventTime4, "proba1"},
	}
	horarium = Horarium{horariumEvents, "la"}

	events = []types.Event{
		{gorm.Model{ID: 0, CreatedAt: time.Time{}, UpdatedAt: time.Time{}}, 0, eventTime1.ToTime(locale), eventTime2.ToTime(locale), []types.LocatedString{{Value: "test0", Language: "de"}, {Value: "proba0", Language: "la"}}},
		{gorm.Model{ID: 1, CreatedAt: time.Time{}, UpdatedAt: time.Time{}}, 1, eventTime3.ToTime(locale), eventTime4.ToTime(locale), []types.LocatedString{{Value: "test1", Language: "de"}, {Value: "proba1", Language: "la"}}},
	}
}

//func TestEventTime_ToTime(t *testing.T) {
//	type fields struct {
//		Year   int
//		Month  int
//		Day    int
//		Hour   int
//		Minute int
//	}
//	type args struct {
//		location *time.Location
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   time.Time
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			evTime := EventTime{
//				Year:   tt.fields.Year,
//				Month:  tt.fields.Month,
//				Day:    tt.fields.Day,
//				Hour:   tt.fields.Hour,
//				Minute: tt.fields.Minute,
//			}
//			if got := evTime.ToTime(tt.args.location); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ToTime() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestHorarium_toEventList(t *testing.T) {
	wantedEvents := events
	for i := range wantedEvents {
		wantedEvents[i].Names = wantedEvents[i].Names[1:2]
	}
	tests := []struct {
		name     string
		horarium Horarium
		idOffset int
		want     []types.Event
	}{
		{"empty horarium", Horarium{Events: []WeekViewEvent{}, Language: "la"}, 0, []types.Event{}},
		{"horarium1", horarium, 0, wantedEvents},
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
	var wantedEvents []types.Event
	_ = copier.Copy(&wantedEvents, &events)
	for i := range wantedEvents {
		wantedEvents[i].ID = i + 10
		wantedEvents[i].Model.ID = uint(i + 10)
	}
	tests := []struct {
		name         string
		dataPath     string
		dataIdOffset int
		want         []types.Event
	}{
		{"test1", "../data/testData/horariumHelper/", 10, wantedEvents},
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
