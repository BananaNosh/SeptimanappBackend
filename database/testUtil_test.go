package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"gorm.io/gorm"
)

// The times of the two events in data/testData/horariumHelper.
// EventTime is a plain value struct, so sharing these between tests is safe.
var (
	eventTime1 = EventTime{Year: 2019, Month: 6, Day: 25, Hour: 16, Minute: 30}
	eventTime2 = EventTime{Year: 2019, Month: 6, Day: 25, Hour: 18, Minute: 0}
	eventTime3 = EventTime{Year: 2019, Month: 6, Day: 25, Hour: 20, Minute: 30}
	eventTime4 = EventTime{Year: 2019, Month: 6, Day: 25, Hour: 21, Minute: 0}
)

// testEvents returns the events of data/testData/horariumHelper with the ids
// ToEventList gives them.
//
// Every call builds new events with new Names: gorm writes primary keys and
// parent ids into the LocatedStrings it stores, and an Event passed by value
// still shares the backing array of its Names. A fixture kept in a package
// variable is therefore mutated by the first test that stores it, and every
// later test sees the mutation.
func testEvents() types.Events {
	return types.Events{
		{
			Model: gorm.Model{ID: 0},
			ID:    0,
			Start: eventTime1.ToTime(util.Locale()),
			End:   eventTime2.ToTime(util.Locale()),
			Names: []types.LocatedString{
				{Value: "test0", Language: "de"},
				{Value: "proba0", Language: "la"},
			},
		},
		{
			Model: gorm.Model{ID: 1},
			ID:    1,
			Start: eventTime3.ToTime(util.Locale()),
			End:   eventTime4.ToTime(util.Locale()),
			Names: []types.LocatedString{
				{Value: "test1", Language: "de"},
				{Value: "proba1", Language: "la"},
			},
		},
	}
}
