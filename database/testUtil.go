package database

import (
	"SeptimanappBackend/types"
	"gorm.io/gorm"
	"time"
)

var eventTime1 EventTime
var eventTime2 EventTime
var eventTime3 EventTime
var eventTime4 EventTime
var events []types.Event

func setupTestEventVariables() {
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
	events = []types.Event{
		{gorm.Model{ID: 0, CreatedAt: time.Time{}, UpdatedAt: time.Time{}}, 0, eventTime1.ToTime(locale), eventTime2.ToTime(locale), []types.LocatedString{{Value: "test0", Language: "de"}, {Value: "proba0", Language: "la"}}},
		{gorm.Model{ID: 1, CreatedAt: time.Time{}, UpdatedAt: time.Time{}}, 1, eventTime3.ToTime(locale), eventTime4.ToTime(locale), []types.LocatedString{{Value: "test1", Language: "de"}, {Value: "proba1", Language: "la"}}},
	}
}
