package main

import (
	"SeptimanappBackend/database"
	Openapi "SeptimanappBackend/openApi"
)

func main() {

	database.InitDatabase()

	Openapi.StartRestApi()
}
