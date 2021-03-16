package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	ID       int `gorm:"primary_key, AUTO_INCREMENT"`
	Start    time.Time
	End      time.Time
	Name     string
	Language string `json:"-"`
}

type Location struct {
	gorm.Model
}

func insertStartEnd(db *gorm.DB) {
	start := time.Date(2021, 7, 31, 16, 30, 0, 0, locale)
	end := time.Date(2021, 8, 7, 14, 0, 0, 0, locale)
	var event Event
	const mainName = "__MAIN__"
	db.FirstOrCreate(&event, Event{Start: start, End: end, Name: mainName})

}

func initDatabase() {
	fmt.Println("Init Database")
	db, err := gorm.Open(sqlite.Open("./data/septimana.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Event{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to auto migrate database")
	}

	insertStartEnd(db)

	insertEventsFromJsonHorarium("./data/", db)
}
