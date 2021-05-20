package database

import (
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
var events []Event

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
	locale, _ = time.LoadLocation("Europe/Berlin")
	horariumEvents := []WeekViewEvent{
		{"e0", eventTime1, eventTime2, "test0"},
		{"ev1", eventTime3, eventTime4, "test1"},
	}
	horarium = Horarium{horariumEvents, "la"}

	events = []Event{
		{gorm.Model{ID: 0, CreatedAt: time.Time{}, UpdatedAt: time.Time{}}, 0, eventTime1.ToTime(locale), eventTime2.ToTime(locale), "test0", "la"},
		{gorm.Model{ID: 1, CreatedAt: time.Time{}, UpdatedAt: time.Time{}}, 1, eventTime3.ToTime(locale), eventTime4.ToTime(locale), "test1", "la"},
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
	tests := []struct {
		name     string
		horarium Horarium
		idOffset int
		want     []Event
	}{
		{"empty horarium", Horarium{Events: []WeekViewEvent{}, Language: "la"}, 0, []Event{}},
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
	var wantedEventsDe []Event
	var wantedEventsLa []Event
	_ = copier.Copy(&wantedEventsDe, &events)
	_ = copier.Copy(&wantedEventsLa, &events)
	for i := range wantedEventsDe {
		wantedEventsDe[i].ID = i + 10
		wantedEventsDe[i].Model.ID = uint(i + 10)
		wantedEventsDe[i].Language = "de"
	}
	for i := range wantedEventsLa {
		wantedEventsLa[i].ID = i + 10 + 2
		wantedEventsLa[i].Model.ID = uint(i + 10 + 2)
	}
	tests := []struct {
		name         string
		dataPath     string
		dataIdOffset int
		want         []Event
	}{
		{"test1", "./testData/horariumHelper/", 10, append(wantedEventsDe, wantedEventsLa...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := EventsFromJsonHoraria(tt.dataPath, tt.dataIdOffset)
			for i := range events {
				events[i].Model.CreatedAt = time.Time{}
				events[i].Model.UpdatedAt = time.Time{}
			}
			if got := events; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventsFromJsonHoraria() = %v, want %v", got, tt.want)
			}
		})
	}
}
