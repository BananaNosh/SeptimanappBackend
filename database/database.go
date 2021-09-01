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

type Repository struct {
	db *gorm.DB
}

func GetRepository() (*Repository, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func insertStartEnd(db *gorm.DB) {
	start := time.Date(2021, 7, 31, 16, 30, 0, 0, util.Locale())
	end := time.Date(2021, 8, 7, 14, 0, 0, 0, util.Locale())
	var event types.Event
	db.FirstOrCreate(&event, types.Event{Start: start, End: end, Names: nil})

}

func (rep *Repository) InitDatabase() {
	fmt.Println("Init Database")
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second,   // Slow SQL threshold
	//		LogLevel:      logger.Info, // Log level
	//	},
	//)
	db := rep.db

	// Migrate the schema
	err := db.AutoMigrate(&types.Event{}, &types.LocatedString{})
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

	err = db.AutoMigrate(types.ApiKeyInfo{})
}

func openDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path.Join(dataPath, "septimana.db")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		//Logger: newLogger,
	})
	if err != nil {
		return db, errors.New("couldn't open database")
	}
	return db, nil
}

func (rep *Repository) GetEvent(id int) (*types.Event, error) {
	var event types.Event
	err := rep.db.First(&event, id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	var locatedStrings []types.LocatedString
	err = rep.db.Model(&types.LocatedString{}).Where("parent_type = ?", "events").Where("parent_id = ?", id).Find(&locatedStrings).Error
	if err != nil {
		return nil, err
	}
	event.Names = locatedStrings
	return &event, nil
}

func (rep *Repository) GetEvents(year *int) ([]types.Event, error) {
	var events []types.Event
	eventsMap := make(map[int]types.Event, len(events))
	var locatedStrings []types.LocatedString
	if year != nil {
		rep.db.Where("SUBSTR(start, 1, 4) = ?", strconv.Itoa(*year)).Find(&events) // TODO make nicer if possible
	} else {
		rep.db.Find(&events)
	}
	for _, e := range events {
		eventsMap[e.ID] = e
	}
	rep.db.Model(&types.LocatedString{}).Where("parent_type = ?", "events").Find(&locatedStrings)

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

	//rep.db.Model(&types.Event{}).Select("users., emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
	//rep.db.Joins("Names").Find(&events)
	//rep.db.Preload("Names").Find(&events)
	//var pairs []struct{
	//	event *types.Event
	//	locatedString *types.LocatedString
	//}
	//rep.db.Model(&types.Event{}).Joins("left join located_strings on located_strings.parent_id = events.id").Find(&pairs)
	//fmt.Println(pairs)
	return events, nil
}

func (rep *Repository) GetLocation(id string) (*types.Location, error) {
	var location types.Location
	err := rep.db.First(&location, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	var titles []types.LocatedString
	var descriptions []types.LocatedString
	err = rep.db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_title").Where("parent_id = ?", id).Find(&titles).Error
	if err != nil {
		return nil, err
	}
	err = rep.db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_description").Where("parent_id = ?", id).Find(&titles).Error
	if err != nil {
		return nil, err
	}
	location.Titles = titles
	location.Titles = descriptions
	return &location, nil
}

func (rep *Repository) GetLocations(overallLocation *types.OverallLocation) ([]types.Location, error) {
	var locations []types.Location
	locationsMap := make(map[string]types.Location, len(locations))
	var titles []types.LocatedString
	var descriptions []types.LocatedString
	if overallLocation != nil {
		rep.db.Where("overall_location = ?", *overallLocation).Find(&locations)
	} else {
		rep.db.Find(&locations)
	}
	for _, l := range locations {
		locationsMap[l.ID] = l
	}

	rep.db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_title").Find(&titles)
	for _, title := range titles {
		parentID := title.ParentID
		location, ok := locationsMap[parentID]
		if ok {
			location.Titles = append(location.Titles, title)
			locationsMap[parentID] = location
		}
	}

	rep.db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_description").Find(&descriptions)
	for _, description := range descriptions {
		parentID := description.ParentID
		location, ok := locationsMap[parentID]
		if ok {
			location.Descriptions = append(location.Descriptions, description)
			locationsMap[parentID] = location
		}
	}

	locations = nil
	for _, l := range locationsMap {
		locations = append(locations, l)
	}

	//db.Model(&types.Event{}).Select("users., emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
	//db.Joins("Names").Find(&locations)
	//db.Preload("Names").Find(&locations)
	//var pairs []struct{
	//	event *types.Event
	//	locatedString *types.LocatedString
	//}
	//db.Model(&types.Event{}).Joins("left join located_strings on located_strings.parent_id = locations.id").Find(&pairs)
	//fmt.Println(pairs)
	return locations, nil
}

func (rep *Repository) StoreSecurityInfo(info types.ApiKeyInfo) {
	rep.db.Create(info)
}

func (rep *Repository) HasApiKeyInfo(info types.ApiKeyInfo) (bool, error) {
	result := rep.db.First(&info, "api_key_hash = ?", info.ApiKeyHash)
	err := result.Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return !errors.Is(result.Error, gorm.ErrRecordNotFound), err
}
