package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"time"
)

const horariumFilePattern = "horarium_(?P<year>\\d+)_(?P<lang>\\w+).json"

type Horarium struct {
	Events   []WeekViewEvent `json:"events"`
	Language string
}

func (horarium Horarium) ToEventList(idOffset int) []types.Event {
	var events = make([]types.Event, 0)
	for i, event := range horarium.Events {
		id := idOffset + i
		events = append(events, types.Event{
			Model: gorm.Model{ID: uint(id), CreatedAt: time.Now(), UpdatedAt: time.Now()},
			ID:    id,
			Start: event.StartTime.ToTime(util.Locale()),
			End:   event.EndTime.ToTime(util.Locale()),
			Names: []types.LocatedString{{
				Value:    event.Name,
				Language: horarium.Language,
			}},
		})
	}
	return events
}

func horariaToEvents(horariumDe *Horarium, horariumLa *Horarium, idOffset int) (_ []types.Event, err error) {
	if (horariumDe.Language != "de" && horariumDe.Language != "") || (horariumLa.Language != "la" && horariumDe.Language != "") {
		err = errors.New("wrong language")
		return
	}
	eventsDe := horariumDe.ToEventList(idOffset)
	eventsLa := horariumLa.ToEventList(idOffset)

	if len(eventsDe) != len(eventsLa) {
		return nil, errors.New("number of events does not match")
	}

	sort.Slice(eventsDe, func(i, j int) bool {
		return eventsDe[i].Start.Before(eventsDe[j].Start) || (eventsDe[i].Start == eventsDe[j].Start && eventsDe[i].End.Before(eventsDe[j].End))
	})
	sort.Slice(eventsLa, func(i, j int) bool {
		return eventsLa[i].Start.Before(eventsLa[j].Start) || (eventsLa[i].Start == eventsLa[j].Start && eventsLa[i].End.Before(eventsLa[j].End))
	})

	for i := range eventsDe {
		deEvent := eventsDe[i]
		laEvent := eventsLa[i]
		if deEvent.Start != laEvent.Start || deEvent.End != laEvent.End {
			return nil, errors.New("events have not identical times")
		}
		eventsDe[i].Names = append(deEvent.Names, laEvent.Names[0])
	}
	return eventsDe, nil
}

type WeekViewEvent struct {
	ID        string    `json:"mId"`
	StartTime EventTime `json:"mStartTime"`
	EndTime   EventTime `json:"mEndTime"`
	Name      string    `json:"mName"`
}

type EventTime struct {
	Year   int `json:"year"`
	Month  int `json:"month"` // 0-11
	Day    int `json:"dayOfMonth"`
	Hour   int `json:"hourOfDay"`
	Minute int `json:"minute"`
}

func (evTime EventTime) ToTime(location *time.Location) time.Time {
	return time.Date(evTime.Year, time.Month(evTime.Month+1), evTime.Day, evTime.Hour, evTime.Month, 0, 0, location)
}

func EventsFromJsonHoraria(dataPath string, dataIdOffset int) ([]types.Event, error) {
	var allEvents []types.Event

	// Call Readdir to get all files.
	outputDirFiles, _ := ioutil.ReadDir(dataPath)

	// compile regex for HorariaFiles
	reg := regexp.MustCompile(horariumFilePattern)
	offset := dataIdOffset
	var filesForYears = make(map[string]struct {
		deFile string
		laFile string
	})
	for _, file := range outputDirFiles {
		match := reg.FindStringSubmatch(file.Name())
		isHorariumFile := len(match) > 1
		if isHorariumFile {
			year := match[1]
			language := match[2]

			pair := filesForYears[year]
			if language == "la" {
				pair.laFile = file.Name()
			} else if language == "de" {
				pair.deFile = file.Name()
			} else {
				return nil, errors.New("unknown language in event")
			}
			filesForYears[year] = pair
		}
	}
	for _, filePair := range filesForYears {
		horariumDe, fileErr := readHorariumFromFile(path.Join(dataPath, filePair.deFile))
		horariumLa, fileErr2 := readHorariumFromFile(path.Join(dataPath, filePair.laFile))
		if fileErr == nil && fileErr2 == nil {
			horariumDe.Language = "de"
			horariumLa.Language = "la"
			events, parseErr := horariaToEvents(&horariumDe, &horariumLa, offset)

			if parseErr == nil {
				offset += len(events)
				allEvents = append(allEvents, events...)
			}
		}
	}
	return allEvents, nil
}

func readHorariumFromFile(filePath string) (Horarium, error) {
	var horarium Horarium
	if horariumFile, err := os.Open(filePath); err == nil {
		defer horariumFile.Close()

		// create a new decoder
		err = json.NewDecoder(horariumFile).Decode(&horarium)
	}
	return horarium, nil
}
