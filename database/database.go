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
	Db *gorm.DB
}

func GetRepository() (Repository, error) {
	db, err := openDB()
	return Repository{Db: db}, err
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

func insertStartEnd(db *gorm.DB) {
	start := time.Date(2021, 7, 31, 16, 30, 0, 0, util.Locale())
	end := time.Date(2021, 8, 7, 14, 0, 0, 0, util.Locale())
	var event types.Event
	db.FirstOrCreate(&event, types.Event{Start: start, End: end, Names: nil})

}

func (rep Repository) InitDatabase() {
	rep.InitDatabaseFromPath(dataPath)
}

func (rep Repository) InitDatabaseFromPath(dataPath string) {
	fmt.Printf("Init Database from path %s\n", dataPath)
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second,   // Slow SQL threshold
	//		LogLevel:      logger.Info, // Log level
	//	},
	//)

	// Migrate the schema
	err := rep.SetupTables()

	db := rep.Db
	firstId := 1
	insertStartEnd(db)
	horariaIdOffset := firstId + 1

	events, err := EventsFromJsonHoraria(dataPath, horariaIdOffset)
	if err != nil {
		fmt.Println(err)
	} else {
		db.Create(events)
	}

	locations := LocationsFromJsonFiles(dataPath)
	//for _, location := range locations { TODO replace with check if existent
	//	Db.Updates(&location)
	//}
	db.Create(locations)
}

func (rep Repository) SetupTables() error {
	err := rep.Db.AutoMigrate(&types.Event{}, &types.LocatedString{})
	err2 := rep.Db.AutoMigrate(&types.Location{})
	err3 := rep.Db.AutoMigrate(types.ApiKeyInfo{})
	if err != nil || err2 != nil || err3 != nil {
		panic("failed to auto migrate database")
	}
	return err
}

func (rep Repository) GetEvent(id int) (*types.Event, error) {
	var event types.Event
	err := rep.Db.First(&event, id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	var locatedStrings []types.LocatedString
	err = rep.Db.Model(&types.LocatedString{}).Where("parent_type = ?", "events").Where("parent_id = ?", id).Find(&locatedStrings).Error
	if err != nil {
		return nil, err
	}
	event.Names = locatedStrings
	if len(event.Names) == 0 { // event without information - should be a total-septimana event
		return nil, nil
	}
	event.Model = gorm.Model{}
	return &event, nil
}

func (rep Repository) GetEvents(year *int) ([]types.Event, error) {
	var events []types.Event
	eventsMap := make(map[int]types.Event, len(events))
	var locatedStrings []types.LocatedString
	if year != nil {
		rep.Db.Where("SUBSTR(start, 1, 4) = ?", strconv.Itoa(*year)).Find(&events) // TODO make nicer if possible
	} else {
		rep.Db.Find(&events)
	}
	for _, e := range events {
		eventsMap[e.ID] = e
	}
	rep.Db.Model(&types.LocatedString{}).Where("parent_type = ?", "events").Find(&locatedStrings)

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
		if len(e.Names) > 0 {
			events = append(events, e) // event without names is useless or a total-septimana event
		}
	}

	//rep.Db.Model(&types.Event{}).Select("users., emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
	//rep.Db.Joins("Names").Find(&events)
	//rep.Db.Preload("Names").Find(&events)
	//var pairs []struct{
	//	event *types.Event
	//	locatedString *types.LocatedString
	//}
	//rep.Db.Model(&types.Event{}).Joins("left join located_strings on located_strings.parent_id = events.id").Find(&pairs)
	//fmt.Println(pairs)
	return events, nil
}

func (rep Repository) AddEvent(event types.Event) (int, error) {
	ids, err := rep.AddEvents(types.Events{event})
	var id int
	if err == nil {
		id = ids[0]
	}
	return id, err
}

func (rep Repository) AddEvents(events types.Events) ([]int, error) {
	for _, ev := range events {
		if ev.ID != 0 {
			return nil, errors.New("event must not have an ID before adding")
		}
	}
	err := rep.Db.Create(&events).Error
	var ids []int
	for _, ev := range events {
		ids = append(ids, ev.ID)
	}
	return ids, err
}

func (rep Repository) GetLocation(id string) (*types.Location, error) {
	var location types.Location
	err := rep.Db.First(&location, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	var titles []types.LocatedString
	var descriptions []types.LocatedString
	err = rep.Db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_title").Where("parent_id = ?", id).Find(&titles).Error
	if err != nil {
		return nil, err
	}
	err = rep.Db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_description").Where("parent_id = ?", id).Find(&titles).Error
	if err != nil {
		return nil, err
	}
	location.Titles = titles
	location.Titles = descriptions
	return &location, nil
}

func (rep Repository) GetLocations(overallLocation *types.OverallLocation) ([]types.Location, error) {
	var locations []types.Location
	locationsMap := make(map[string]types.Location, len(locations))
	var titles []types.LocatedString
	var descriptions []types.LocatedString
	if overallLocation != nil {
		rep.Db.Where("overall_location = ?", *overallLocation).Find(&locations)
	} else {
		rep.Db.Find(&locations)
	}
	for _, l := range locations {
		locationsMap[l.ID] = l
	}

	rep.Db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_title").Find(&titles)
	for _, title := range titles {
		parentID := title.ParentID
		location, ok := locationsMap[parentID]
		if ok {
			location.Titles = append(location.Titles, title)
			locationsMap[parentID] = location
		}
	}

	rep.Db.Model(&types.LocatedString{}).Where("parent_type = ?", "locations_description").Find(&descriptions)
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

	//Db.Model(&types.Event{}).Select("users., emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
	//Db.Joins("Names").Find(&locations)
	//Db.Preload("Names").Find(&locations)
	//var pairs []struct{
	//	event *types.Event
	//	locatedString *types.LocatedString
	//}
	//Db.Model(&types.Event{}).Joins("left join located_strings on located_strings.parent_id = locations.id").Find(&pairs)
	//fmt.Println(pairs)
	return locations, nil
}

func (rep Repository) StoreSecurityInfo(info types.ApiKeyInfo) {
	rep.Db.Create(info)
}

func (rep Repository) HasApiKeyInfo(info types.ApiKeyInfo) (bool, error) {
	result := rep.Db.First(&info, "api_key_hash = ?", info.ApiKeyHash)
	err := result.Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return !errors.Is(result.Error, gorm.ErrRecordNotFound), err
}
