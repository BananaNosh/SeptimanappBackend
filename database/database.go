package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
	"strconv"
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
	db, err := openDB()
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

func openDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path.Join(dataPath, "septimana.db")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		//Logger: newLogger,
	})
}

func GetEvent(id int) (*types.Event, error) {
	db, err := openDB()
	if err != nil {
		return nil, errors.New("couldn't open database")
	}

	var event types.Event
	err = db.Find(&event, id).Error
	if err != nil {
		fmt.Print(err)
	}

	var locatedStrings []types.LocatedString
	err = db.Model(&types.LocatedString{}).Where("parent_type = ?", "events").Where("parent_id = ?", id).Find(&locatedStrings).Error
	if err != nil {
		fmt.Print(err)
	}
	event.Names = locatedStrings
	return &event, nil
}

func GetEvents(year *int) ([]types.Event, error) {
	db, err := openDB()
	if err != nil {
		return nil, errors.New("couldn't open database")
	}

	var events []types.Event
	eventsMap := make(map[int]types.Event, len(events))
	var locatedStrings []types.LocatedString
	if year != nil {
		db.Where("SUBSTR(start, 1, 4) = ?", strconv.Itoa(*year)).Find(&events) // TODO make nicer if possible
	} else {
		db.Find(&events)
	}
	for _, e := range events {
		eventsMap[e.ID] = e
	}
	db.Model(&types.LocatedString{}).Where("parent_type = ?", "events").Find(&locatedStrings)

	for _, locString := range locatedStrings {
		parentID, err := strconv.Atoi(locString.ParentID)
		if err == nil {
			event, ok := eventsMap[parentID]
			if ok {
				event.Names = append(event.Names, locString)
				eventsMap[parentID] = event
			}
		}
	}
	events = nil
	for _, e := range eventsMap {
		events = append(events, e)
	}

	//db.Model(&types.Event{}).Select("users., emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
	//db.Joins("Names").Find(&events)
	//db.Preload("Names").Find(&events)
	//var pairs []struct{
	//	event *types.Event
	//	locatedString *types.LocatedString
	//}
	//db.Model(&types.Event{}).Joins("left join located_strings on located_strings.parent_id = events.id").Find(&pairs)
	//fmt.Println(pairs)
	return events, nil
}
