package main

import (
	"encoding/json"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

const locationFilePattern = "locations_(?P<location>(\\w|\\d)+)\\.json"

func LocationsFromJsonFiles(dataPath string) []Location {
	var allLocations []Location

	// Call Readdir to get all files.
	outputDirFiles, _ := ioutil.ReadDir(dataPath)

	// compile regex for HorariaFiles
	reg := regexp.MustCompile(locationFilePattern)
	for _, file := range outputDirFiles {
		match := reg.FindStringSubmatch(file.Name())
		isLocationFile := len(match) > 1
		if isLocationFile {
			overallLocation := match[1]

			// open the file pointer
			filePath := path.Join(dataPath, file.Name())

			if locations, err := readLocationsFromFile(filePath); err == nil {
				// set correct overallLocation
				for i := range locations {
					locations[i].OverallLocation = overallLocation
				}
				allLocations = append(allLocations, locations...)
			}
		}
	}
	return allLocations
}

func readLocationsFromFile(filePath string) (_ []Location, err error) {
	var locations = struct {
		Locations []Location
	}{}
	if locationsFile, err := os.Open(filePath); err == nil {
		defer locationsFile.Close()

		// create a new decoder
		err = json.NewDecoder(locationsFile).Decode(&locations)
	}
	return locations.Locations, nil
}

func (location *Location) UnmarshalJSON(data []byte) (err error) {
	/**
	Unmarshal json bytes to location
	*/
	var auxiliaryLocation struct {
		Id             string
		TitleMap       map[string]string
		DescriptionMap map[string]string
		IsMain         bool
		Coordinates    map[string]float32
	}
	if err = json.Unmarshal(data, &auxiliaryLocation); err == nil {
		location.ID = auxiliaryLocation.Id
		location.Longitude = auxiliaryLocation.Coordinates["mLongitude"]
		location.Latitude = auxiliaryLocation.Coordinates["mLatitude"]
		location.Altitude = auxiliaryLocation.Coordinates["mAltitude"]
		location.IsMain = auxiliaryLocation.IsMain
		location.Titles = locationStringsFromMap(auxiliaryLocation.TitleMap)
		location.Descriptions = locationStringsFromMap(auxiliaryLocation.DescriptionMap)
	}

	return err
}

func locationStringsFromMap(stringMap map[string]string) []LocationString {
	var locationStrings []LocationString
	for k, v := range stringMap {
		locationStrings = append(locationStrings, LocationString{
			Model:    gorm.Model{},
			Value:    v,
			Language: k,
		})
	}
	return locationStrings
}
