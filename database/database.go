package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
	"time"
)

const dataPath = "./data"

func insertStartEnd(db *gorm.DB) {
	start := time.Date(2021, 7, 31, 16, 30, 0, 0, util.Locale())
	end := time.Date(2021, 8, 7, 14, 0, 0, 0, util.Locale())
	var event types.Event
	db.FirstOrCreate(&event, types.Event{Start: start, End: end, Names: nil})

}

func InitDatabase() {
	fmt.Println("Init Database")
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second,   // Slow SQL threshold
	//		LogLevel:      logger.Info, // Log level
	//	},
	//)
	db, err := gorm.Open(sqlite.Open(path.Join(dataPath, "septimana.db")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		//Logger: newLogger,
	})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&types.Event{}, &types.LocatedString{})
	if err != nil {
		panic("failed to auto migrate database")
	}

	firstId := 1
	insertStartEnd(db)
	horariaIdOffset := firstId + 1

	events, err := EventsFromJsonHoraria(dataPath, horariaIdOffset)
	if err != nil {
		fmt.Println(err)
	} else {
		db.Create(events)
	}

	err = db.AutoMigrate(&types.Location{})
	if err != nil {
		panic("failed to auto migrate database")
	}

	locations := LocationsFromJsonFiles(dataPath)
	//for _, location := range locations { TODO replace with check if existent
	//	db.Updates(&location)
	//}
	db.Create(locations)

}
