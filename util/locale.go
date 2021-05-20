package util

import "time"

var locale *time.Location

func Locale() *time.Location {
	if locale == nil {
		location, err := time.LoadLocation("Europe/Berlin")
		if err != nil {
			location = time.Local
		}
		locale = location
	}
	return locale
}
