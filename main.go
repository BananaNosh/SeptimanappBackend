package main

import (
	"fmt"
	"time"
)

var locale *time.Location

func main() {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		location = time.Local
	}
	locale = location
	fmt.Println("Test")
	initDatabase()
}
