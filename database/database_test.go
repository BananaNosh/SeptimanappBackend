package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"github.com/jinzhu/copier"
	"testing"
	"time"
)

var wantedEvents types.Events

func databaseTestSetup() {
	setupTestEventVariables()
	_ = copier.Copy(&wantedEvents, &events)
	wantedEvents[0].ID = 2
	wantedEvents[1].ID = 3
	ev3 := types.Event{
		Start: eventTime1.ToTime(util.Locale()),
		End:   eventTime2.ToTime(util.Locale()),
		Names: []types.LocatedString{
			{Value: "test3", Language: "de"},
			{Value: "proba3", Language: "lat"},
		},
	}
	ev4 := types.Event{
		ID:    5,
		Start: eventTime1.ToTime(util.Locale()),
		End:   eventTime2.ToTime(util.Locale()),
		Names: []types.LocatedString{
			{Value: "test4", Language: "de"},
			{Value: "proba4", Language: "lat"},
		},
	}
	ev5 := types.Event{
		Start: eventTime1.ToTime(util.Locale()).Add(time.Duration(24*365) * time.Hour),
		End:   eventTime2.ToTime(util.Locale()).Add(time.Duration(24*365) * time.Hour),
		Names: []types.LocatedString{
			{Value: "test4", Language: "de"},
			{Value: "proba4", Language: "lat"},
		},
	}
	wantedEvents = append(wantedEvents, ev3, ev4, ev5)
}

func setupDatabaseMock(t *testing.T) Repository {
	repository, err := GetMockRepository()
	if err != nil {
		t.Error(err)
	}
	InitMockRepository(repository)
	return repository
}

func TestRepository_GetEvent(t *testing.T) {
	databaseTestSetup()
	tests := []struct {
		name    string
		eventId int
		want    *types.Event
		wantErr bool
	}{
		{"unknown id", 4, nil, false},
		{"totalEvent", 1, nil, false}, // Do not return TotalSeptimana event
		{"e0", 2, &wantedEvents[0], false},
		{"ev1", 3, &wantedEvents[1], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			got, err := rep.GetEvent(tt.eventId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !util.EqualEvent(got, tt.want, true) {
				t.Errorf("\nGetEvent() got = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetEvents(t *testing.T) {
	databaseTestSetup()
	year1 := 1900
	year2 := 2020
	tests := []struct {
		name    string
		year    *int
		want    []types.Event
		wantErr bool
	}{
		{"all", nil, wantedEvents[:2], false},
		{"old", &year1, types.Events{}, false},
		{"current", &year2, types.Events{wantedEvents[4]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			if tt.name == "current" {
				_, _ = rep.AddEvent(wantedEvents[4])
			}
			got, err := rep.GetEvents(tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !util.EqualEvents(got, tt.want, tt.name != "current") {
				t.Errorf("GetEvents() got = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_AddEvent(t *testing.T) {
	databaseTestSetup()
	tests := []struct {
		name    string
		event   types.Event
		wantErr bool
	}{
		{"add ev3", wantedEvents[2], false},
		{"add ev4 with ID", wantedEvents[3], true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			id, err := rep.AddEvent(tt.event)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AddEvents() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			newEv, err := rep.GetEvent(id)
			if err != nil {
				t.Errorf("AddEvents() error = %v", err)
				return
			}
			if !util.EqualEvent(newEv, &tt.event, false) {
				t.Errorf("AddEvents() event = %v,\n originalEvent %v", newEv, tt.event)
				return
			}
		})
	}
}

func TestRepository_AddEvents(t *testing.T) {
	databaseTestSetup()
	tests := []struct {
		name                   string
		events                 types.Events
		wantErr                bool
		wantedTotalEventsCount int
	}{
		{"one wrong", wantedEvents[2:4], true, 2},
		{"both existing", wantedEvents[:2], true, 2},
		{"add two correct", types.Events{wantedEvents[2], wantedEvents[4]}, false, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			ids, err := rep.AddEvents(tt.events)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AddEvents() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			total, err := rep.GetEvents(nil)
			if err != nil {
				t.Errorf("AddEvents() -> could not get all events")
			}
			if len(total) != tt.wantedTotalEventsCount {
				t.Errorf("AddEvents() totalEvents = %v, wantedTotalEventsCount %v", total, tt.wantedTotalEventsCount)
			}
			for i, id := range ids {
				event, err := rep.GetEvent(id)
				if err != nil || !util.EqualEvent(event, &tt.events[i], false) {
					t.Errorf("AddEvents() event = %v, originalEvent %v, err = %v", event, tt.events[i], err)
				}
			}
		})
	}
}

func TestRepository_DeleteEvent(t *testing.T) {
	databaseTestSetup()
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"delete ev2", 2, false},
		{"delete ev10", 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			existingEvent, _ := rep.GetEvent(tt.id)
			var existingNames []types.LocatedString
			if existingEvent != nil {
				existingNames = existingEvent.Names
			}
			if err := rep.DeleteEvent(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			event, err := rep.GetEvent(tt.id)
			if event != nil {
				t.Errorf("DeleteEvent() got deleted event = %v, error = %v", event, err)
			}
			for _, exName := range existingNames {
				var name types.LocatedString
				err := rep.Db.First(&name, exName.ID).Error
				if err == nil {
					t.Errorf("DeleteEvent() got deleted name = %v", name)
				}
			}
		})
	}
}

func TestRepository_DeleteEvents(t *testing.T) {
	databaseTestSetup()
	tests := []struct {
		name    string
		ids     []int
		wantErr bool
	}{
		{"delete ev2 and ev3", []int{2, 3}, false},
		{"delete ev3 and 10", []int{3, 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			if err := rep.DeleteEvents(tt.ids); (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, id := range tt.ids {
				event, err := rep.GetEvent(id)
				if event != nil {
					t.Errorf("DeleteEvent() got deleted event = %v, error = %v", event, err)
				}
			}
		})
	}
}

func TestRepository_UpdateEvent(t *testing.T) {
	databaseTestSetup()
	newEvent2 := wantedEvents[0]
	newEvent2.Start = eventTime3.ToTime(util.Locale())
	newEvent2.End = eventTime4.ToTime(util.Locale())
	newEvent2.Names = append(newEvent2.Names, types.LocatedString{Value: "huh", Language: "kling"})
	tests := []struct {
		name    string
		event   types.Event
		wantErr bool
	}{
		{"update ev2", newEvent2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			if err := rep.UpdateEvent(tt.event); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			event, _ := rep.GetEvent(tt.event.ID)
			if !util.EqualEvent(event, &newEvent2, true) {
				t.Errorf("UpdateEvent() wanted = %v, got event = %v", newEvent2, event)
			}
		})
	}
}
