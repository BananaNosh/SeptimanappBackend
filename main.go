package main

import (
	"SeptimanappBackend/database"
	Openapi "SeptimanappBackend/openApi"
	"fmt"
)

func main() {

	repository, err := database.GetRepository()
	if err != nil {
		fmt.Println(err)
	}
	repository.InitDatabase()

	Openapi.StartRestApi()
}
