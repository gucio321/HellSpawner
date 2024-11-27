package main

import (
	"log"

	"github.com/gucio321/HellSpawner/pkg/app"
)

func main() {
	log.SetFlags(log.Lshortfile)

	app, err := app.Create()
	if err != nil {
		log.Fatal(err)
	} else if app == nil {
		return // we've terminated early
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
