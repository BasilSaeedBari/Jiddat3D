package main

import (
	"log"
	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}
	
	record, err := app.FindFirstRecordByFilter("content", "slug = 'what-is-reprap-pakistan'")
	if err != nil {
		log.Fatal("Not found:", err)
	}
	
	log.Printf("ID: %v, Type: %v, Published: %v\n", record.Id, record.Get("type"), record.Get("published"))
}
