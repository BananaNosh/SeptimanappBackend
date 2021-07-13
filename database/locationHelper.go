package database

import (
	"SeptimanappBackend/types"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

const locationFilePattern = "locations_(?P<location>(\\w|\\d)+)\\.json"

func LocationsFromJsonFiles(dataPath string) []types.Location {
	var allLocations []types.Location

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
					locations[i].OverallLocation = types.OverallLocation(overallLocation)
				}
				allLocations = append(allLocations, locations...)
			}
		}
	}
	return allLocations
}

func readLocationsFromFile(filePath string) (_ []types.Location, err error) {
	var locations = struct {
		Locations []types.Location
	}{}
	if locationsFile, err := os.Open(filePath); err == nil {
		defer locationsFile.Close()

		// create a new decoder
		err = json.NewDecoder(locationsFile).Decode(&locations)
	}
	return locations.Locations, nil
}
