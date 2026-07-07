package main

import (
	"log"
	"jiddat3d/internal/app"
)

func main() {
	pb := app.NewPocketBaseApp()
	
	// Just boot the DB and update the record
	if err := pb.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	record, err := pb.FindFirstRecordByFilter("site_settings", "1=1")
	if err != nil {
		log.Fatal(err)
	}

	record.Set("whatsapp_business_number", "+923001318112")
	record.Set("discord_url", "https://discord.gg/rvsaqF2Q7p")
	
	// Does whatsapp_community_url exist?
	// If not, we will just use it in the template manually.
	
	if err := pb.Save(record); err != nil {
		log.Fatal(err)
	}
	
	log.Println("Settings updated!")
}
