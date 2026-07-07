package main

import (
	"log"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/dbx"
)

func main() {
	app := pocketbase.New()
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}
	
	records, err := app.FindRecordsByFilter("content", "type = {:type} && published = true", "-published_at", 100, 0, dbx.Params{"type": "blog"})
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	
	log.Printf("Found %d records with published = true\n", len(records))

	records2, err2 := app.FindRecordsByFilter("content", "type = {:type}", "-created", 100, 0, dbx.Params{"type": "blog"})
	if err2 != nil {
		log.Printf("Error2: %v\n", err2)
	}
	
	log.Printf("Found %d records WITHOUT published filter\n", len(records2))
}
