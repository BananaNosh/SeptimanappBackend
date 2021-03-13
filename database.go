package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	Id    int `gorm:"primary_key, AUTO_INCREMENT"`
	Start time.Time
	End   time.Time
	Name  string
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

	//event := Event{Name: "Test", Start: time.Now(), End: time.Date(2010, 10, 01, 10, 10, 10, 10, time.Local)}
	//fmt.Println(event)
	//db.Create(&event)
}
