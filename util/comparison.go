package util

import (
	"SeptimanappBackend/types"
	"sort"
)

func EqualEvents(events1 types.Events, events2 types.Events, checkIds bool) bool {
	if len(events1) != len(events2) {
		return false
	}
	sort.Slice(events1, func(i, j int) bool {
		return events1[i].ID < events1[j].ID
	})
	sort.Slice(events2, func(i, j int) bool {
		return events2[i].ID < events2[j].ID
	})
	for i := range events1 {
		e1 := events1[i]
		e2 := events2[i]
		if !EqualEvent(&e1, &e2, checkIds) {
			return false
		}
	}
	return true
}

func EqualEvent(event1 *types.Event, event2 *types.Event, checkId bool) bool {
	if event1 == nil || event2 == nil {
		return event1 == event2
	}
	if !event1.Start.Equal(event2.Start) || !event1.End.Equal(event2.End) {
		return false
	}
	return (!checkId || event1.ID == event2.ID) && EqualLocatedStrings(event1.Names, event2.Names)
}

// EqualLocatedStrings compares located strings as a language to value mapping,
// which is also how they are sent over the api.
//
// Their order carries no meaning and cannot be normalised by sorting on the id:
// strings decoded from json or built as a fixture have no id yet, so sorting
// leaves them in whatever order they were built in - for json that is the
// iteration order of a map, which differs from run to run.
func EqualLocatedStrings(strings1 []types.LocatedString, strings2 []types.LocatedString) bool {
	if len(strings2) != len(strings1) {
		return false
	}
	byLanguage := make(map[string]string, len(strings1))
	for _, s := range strings1 {
		byLanguage[s.Language] = s.Value
	}
	for _, s := range strings2 {
		if value, found := byLanguage[s.Language]; !found || value != s.Value {
			return false
		}
	}
	return true
}

func EqualLocatedString(s1 *types.LocatedString, s2 *types.LocatedString) bool {
	if s1 == nil || s2 == nil {
		return s1 == s2
	}
	return s1.Language == s2.Language && s1.Value == s2.Value
}

func EqualLocations(locations1 []types.Location, locations2 []types.Location) bool {
	if len(locations1) != len(locations2) {
		return false
	}
	sort.Slice(locations1, func(i, j int) bool {
		return locations1[i].ID < locations1[j].ID
	})
	sort.Slice(locations2, func(i, j int) bool {
		return locations2[i].ID < locations2[j].ID
	})
	for i := range locations1 {
		if !EqualLocation(&locations1[i], &locations2[i]) {
			return false
		}
	}
	return true
}

// EqualLocation compares the contents of two locations, ignoring gorm's
// bookkeeping - GetLocation, unlike GetEvent, does not blank it.
func EqualLocation(location1 *types.Location, location2 *types.Location) bool {
	if location1 == nil || location2 == nil {
		return location1 == location2
	}
	return location1.ID == location2.ID &&
		location1.OverallLocation == location2.OverallLocation &&
		location1.Longitude == location2.Longitude &&
		location1.Latitude == location2.Latitude &&
		location1.Altitude == location2.Altitude &&
		location1.IsMain == location2.IsMain &&
		EqualLocatedStrings(location1.Titles, location2.Titles) &&
		EqualLocatedStrings(location1.Descriptions, location2.Descriptions)
}
