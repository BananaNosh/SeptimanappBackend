package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"testing"
	"time"
)

// Indices into wantedTestEvents().
const (
	seededEvent0     = iota // stored by InitMockRepository, gets id 2
	seededEvent1            // stored by InitMockRepository, gets id 3
	newEvent                // not stored yet
	newEventWithID          // not stored yet, but carries an id - adding it must fail
	newEventNextYear        // not stored yet, one year after the seeded events
)

// wantedTestEvents returns the two events InitMockRepository seeds - with the ids
// the database gives them, id 1 belonging to the total-septimana event - followed
// by events that are not stored yet. Fresh on every call, see testEvents.
func wantedTestEvents() types.Events {
	wanted := testEvents()
	wanted[seededEvent0].ID = 2
	wanted[seededEvent1].ID = 3
	return append(wanted,
		types.Event{
			Start: eventTime1.ToTime(util.Locale()),
			End:   eventTime2.ToTime(util.Locale()),
			Names: []types.LocatedString{
				{Value: "test3", Language: "de"},
				{Value: "proba3", Language: "lat"},
			},
		},
		types.Event{
			ID:    5,
			Start: eventTime1.ToTime(util.Locale()),
			End:   eventTime2.ToTime(util.Locale()),
			Names: []types.LocatedString{
				{Value: "test4", Language: "de"},
				{Value: "proba4", Language: "lat"},
			},
		},
		types.Event{
			Start: eventTime1.ToTime(util.Locale()).Add(time.Duration(24*365) * time.Hour),
			End:   eventTime2.ToTime(util.Locale()).Add(time.Duration(24*365) * time.Hour),
			Names: []types.LocatedString{
				{Value: "test4", Language: "de"},
				{Value: "proba4", Language: "lat"},
			},
		},
	)
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
	tests := []struct {
		name      string
		eventId   int
		wantIndex int // -1 if no event is expected
		wantErr   bool
	}{
		{"unknown id", 4, -1, true},
		{"totalEvent", 1, -1, true}, // Do not return TotalSeptimana event
		{"e0", 2, seededEvent0, false},
		{"ev1", 3, seededEvent1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var want *types.Event
			if tt.wantIndex >= 0 {
				wanted := wantedTestEvents()
				want = &wanted[tt.wantIndex]
			}
			rep := setupDatabaseMock(t)
			got, err := rep.GetEvent(tt.eventId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !util.EqualEvent(got, want, true) {
				t.Errorf("\nGetEvent() got = %v,\n want %v", got, want)
			}
		})
	}
}

func TestRepository_GetEvents(t *testing.T) {
	year1 := 1900
	year2 := 2020
	tests := []struct {
		name        string
		year        *int
		wantIndices []int
		addIndex    int // event to store before the call, -1 for none
		wantErr     bool
	}{
		{"all", nil, []int{seededEvent0, seededEvent1}, -1, false},
		{"old", &year1, []int{}, -1, false},
		{"current", &year2, []int{newEventNextYear}, newEventNextYear, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			if tt.addIndex >= 0 {
				// its own fixture: storing an event mutates the Names it holds
				_, _ = rep.AddEvent(wantedTestEvents()[tt.addIndex])
			}
			wanted := wantedTestEvents()
			want := types.Events{}
			for _, i := range tt.wantIndices {
				want = append(want, wanted[i])
			}

			got, err := rep.GetEvents(tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// the stored event got its id from the database, so it cannot be checked
			if !util.EqualEvents(got, want, tt.addIndex < 0) {
				t.Errorf("GetEvents() got = %v,\n want %v", got, want)
			}
		})
	}
}

func TestRepository_AddEvent(t *testing.T) {
	tests := []struct {
		name       string
		eventIndex int
		wantErr    bool
	}{
		{"add ev3", newEvent, false},
		{"add ev4 with ID", newEventWithID, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := wantedTestEvents()[tt.eventIndex]
			rep := setupDatabaseMock(t)
			id, err := rep.AddEvent(event)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AddEvent() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			newEv, err := rep.GetEvent(id)
			if err != nil {
				t.Errorf("AddEvent() error = %v", err)
				return
			}
			if !util.EqualEvent(newEv, &event, false) {
				t.Errorf("AddEvent() event = %v,\n originalEvent %v", newEv, event)
			}
		})
	}
}

func TestRepository_AddEvents(t *testing.T) {
	tests := []struct {
		name                   string
		eventIndices           []int
		wantErr                bool
		wantedTotalEventsCount int
	}{
		{"one wrong", []int{newEvent, newEventWithID}, true, 2},
		{"both existing", []int{seededEvent0, seededEvent1}, true, 2},
		{"add two correct", []int{newEvent, newEventNextYear}, false, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wanted := wantedTestEvents()
			events := types.Events{}
			for _, i := range tt.eventIndices {
				events = append(events, wanted[i])
			}
			rep := setupDatabaseMock(t)
			ids, err := rep.AddEvents(events)
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
				t.Errorf("AddEvents() totalCount %v, wantedTotalEventsCount %v\nevents = %v", len(total), tt.wantedTotalEventsCount, total)
			}
			for i, id := range ids {
				event, err := rep.GetEvent(id)
				if err != nil || !util.EqualEvent(event, &events[i], false) {
					t.Errorf("AddEvents() event = %v, originalEvent %v, err = %v", event, events[i], err)
				}
			}
		})
	}
}

func TestRepository_GetLocation(t *testing.T) {
	tests := []struct {
		name       string
		locationId string
		wantIndex  int // -1 if no location is expected
		wantErr    bool
	}{
		{"unknown id", "nothing", -1, true},
		{"test1", "test1", 0, false},
		{"test2", "test2", 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var want *types.Location
			if tt.wantIndex >= 0 {
				wanted := testLocations()
				want = &wanted[tt.wantIndex]
			}
			rep := setupDatabaseMock(t)
			got, err := rep.GetLocation(tt.locationId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !util.EqualLocation(got, want) {
				t.Errorf("\nGetLocation() got = %v,\n want %v", got, want)
			}
		})
	}
}

func TestRepository_GetLocations(t *testing.T) {
	over1 := types.OverallLocation("over1")
	unknown := types.OverallLocation("nothing")
	tests := []struct {
		name            string
		overallLocation *types.OverallLocation
		wantIndices     []int
	}{
		{"all", nil, []int{0, 1}},
		{"of over1", &over1, []int{0, 1}},
		{"of unknown overallLocation", &unknown, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wanted := testLocations()
			want := []types.Location{}
			for _, i := range tt.wantIndices {
				want = append(want, wanted[i])
			}
			rep := setupDatabaseMock(t)
			got, err := rep.GetLocations(tt.overallLocation)
			if err != nil {
				t.Errorf("GetLocations() error = %v", err)
				return
			}
			if !util.EqualLocations(got, want) {
				t.Errorf("GetLocations() got = %v,\n want %v", got, want)
			}
		})
	}
}

func TestRepository_DeleteEvent(t *testing.T) {
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
	// moves the stored event 2 to another time and gives it a third name
	updatedEvent2 := func() types.Event {
		event := wantedTestEvents()[seededEvent0]
		event.Start = eventTime3.ToTime(util.Locale())
		event.End = eventTime4.ToTime(util.Locale())
		event.Names = append(event.Names, types.LocatedString{Value: "huh", Language: "kling"})
		return event
	}
	unknownEvent10 := func() types.Event {
		return types.Event{ID: 10, Start: time.Time{}, End: time.Time{}, Names: nil}
	}
	tests := []struct {
		name    string
		event   func() types.Event
		wantErr bool
	}{
		{"update ev2", updatedEvent2, false},
		{"update ev10", unknownEvent10, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.event()
			rep := setupDatabaseMock(t)
			err := rep.UpdateEvent(event)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			got, _ := rep.GetEvent(event.ID)
			if !util.EqualEvent(got, &event, true) {
				t.Errorf("UpdateEvent() wanted = %v, got event = %v", event, got)
			}
		})
	}
}
