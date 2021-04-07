package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
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
	ID              string `gorm:"primary_key"`
	OverallLocation string
	Longitude       float32
	Latitude        float32
	Altitude        float32
	IsMain          bool
	Titles          []LocationString `gorm:"foreignKey:LocationID"`
	Descriptions    []LocationString `gorm:"foreignKey:LocationID"`
}

type LocationString struct {
	gorm.Model
	ID         int `gorm:"primary_key, AUTO_INCREMENT"`
	Value      string
	Language   string
	LocationID string
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
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second,   // Slow SQL threshold
	//		LogLevel:      logger.Info, // Log level
	//	},
	//)
	db, err := gorm.Open(sqlite.Open(path.Join(dataPath, "septimana.db")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // TODO needed ? or even bad
		//Logger: newLogger,
	})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Event{})
	if err != nil {
		panic("failed to auto migrate database")
	}

	firstId := 1
	insertStartEnd(db)
	horariaIdOffset := firstId + 1

	events := EventsFromJsonHoraria(dataPath, horariaIdOffset)
	db.Create(events)

	err = db.AutoMigrate(&Location{}, &LocationString{})
	if err != nil {
		panic("failed to auto migrate database")
	}

	locations := LocationsFromJsonFiles(dataPath)
	fmt.Println(locations)
}
