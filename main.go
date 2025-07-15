package main

import (
	"ari-inventory/api"
	"ari-inventory/database"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
	log.Println("start")

	err := database.NewDatabase("sqlite3", "service_discovery.db")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	http.ListenAndServe(":9190", mux)

}

