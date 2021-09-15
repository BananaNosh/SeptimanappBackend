package util

import (
	"SeptimanappBackend/types"
)

func EqualEvents(events1 types.Events, events2 types.Events) bool {
	if len(events1) != len(events2) {
		return false
	}
	for i := range events1 {
		e1 := events1[i]
		e2 := events2[i]
		if !EqualEvent(&e1, &e2) {
			return false
		}
	}
	return true
}

func EqualEvent(event1 *types.Event, event2 *types.Event) bool {
	if event1 == nil || event2 == nil {
		return event1 == event2
	}
	if !event1.Start.Equal(event2.Start) || !event1.End.Equal(event2.End) {
		return false
	}
	return event1.ID == event2.ID && EqualLocatedStrings(event1.Names, event2.Names)
}

func EqualLocatedStrings(strings1 []types.LocatedString, strings2 []types.LocatedString) bool {
	if len(strings2) != len(strings1) {
		return false
	}
	for i := range strings1 {
		s1 := strings1[i]
		s2 := strings2[i]
		if !EqualLocatedString(&s1, &s2) {
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
