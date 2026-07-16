package main

import (
	"SeptimanappBackend/database"
	Openapi "SeptimanappBackend/openApi"
	"SeptimanappBackend/security"
	"flag"
	"fmt"
	"os"
)

func main() {
	newKey := flag.Bool("newkey", false, "store a new api key and print it, instead of starting the server")
	flag.Parse()

	repository, err := database.GetRepository()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	repository.InitDatabase()

	if *newKey {
		key, err := security.StoreNewApiKey(repository)
		if err != nil {
			fmt.Printf("could not store a new api key: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("New api key: %s\n", key)
		return
	}

	Openapi.StartRestApi()
}
