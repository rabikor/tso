package main

import (
	"log"
	"treatment-scheme-organizer/internal/rdb"
)

func main() {
	if err := rdb.Migrate(); err != nil {
		log.Panic("failed to migrate", err)
	}

	log.Println("models were migrated")
}
