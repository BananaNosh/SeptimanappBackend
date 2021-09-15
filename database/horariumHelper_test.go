package database

import (
	"SeptimanappBackend/types"
	"github.com/jinzhu/copier"
	"reflect"
	"testing"
	"time"
)

var horarium Horarium

func horariumHelperTestSetup() {
	setupTestEventVariables()

	horariumEvents := []WeekViewEvent{
		{"e0", eventTime1, eventTime2, "proba0"},
		{"ev1", eventTime3, eventTime4, "proba1"},
	}
	horarium = Horarium{horariumEvents, "la"}
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
	horariumHelperTestSetup()
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
	horariumHelperTestSetup()
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
