package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	c "github.com/levelabs/level-go/collection"
	s "github.com/levelabs/level-go/scheduler"
	"log"
	"net/http"
	"time"
)

const (
	collectionBAYC      = "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"
	colectionPartyDegen = "0x4be3223f8708ca6b30d1e8b8926cf281ec83e770"
)

type App struct {
	scheduler *gocron.Scheduler

	collectionQueue *c.CollectionQueue
}

func NewApp(assets map[string]int64) *App {
	scheduler := s.NewScheduler()
	collectionQueue := c.NewCollectionQueue(assets)

	app := App{
		scheduler:       scheduler,
		collectionQueue: collectionQueue,
	}

	return &app
}

func main() {
	assets := map[string]int64{
		"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d": time.Now().UnixNano(),
		"0x4be3223f8708ca6b30d1e8b8926cf281ec83e770": time.Now().UnixNano(),
		"0x8a90cab2b38dba80c64b7734e58ee1db38b8992e": time.Now().UnixNano(),
	}

	app := NewApp(assets)

	app.scheduler.Every(5).Seconds().Do(func() {
		fmt.Println("Scheduler starting...")
	})

	s.Start(app.scheduler)

	// c := collection.NewAsset(colectionPartyDegen)
	// c.SetBaseURI()
	// c.QueryAttributes()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
