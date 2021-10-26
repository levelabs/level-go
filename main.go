package main

import (
	"encoding/json"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/ristretto"
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
	scheduler       *gocron.Scheduler
	collectionQueue *c.CollectionQueue

	cache *ristretto.Cache
	db    *badger.DB
}

func NewApp(assets map[string]int64) *App {
	scheduler := s.NewScheduler()
	collectionQueue := c.NewCollectionQueue(assets)

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		log.Fatal(err)
	}

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := App{
		scheduler:       scheduler,
		collectionQueue: collectionQueue,
		cache:           cache,
		db:              db,
	}

	return &app
}

func (app *App) Schedule() {
	app.scheduler.Every(5).Seconds().Do(func() {
		err, asset := app.collectionQueue.RunSequence()
		if err != nil {
			// do something
			fmt.Println("can't run")
			return
		}

		fmt.Println("asset sequenced", asset.Address())

		serialized, err := json.Marshal(asset)
		if err != nil {
			// do something
		}

		address := asset.Address()
		key := []byte(address)
		txn := app.db.NewTransaction(true)

		if err := txn.Set(key, []byte(serialized)); err == badger.ErrTxnTooBig {
			if err := txn.Commit(); err != nil {
				// do something
			}
			txn = app.db.NewTransaction(true)
			if err := txn.Set(key, []byte(serialized)); err != nil {
				// do something
			}
		}
	})

	s.Start(app.scheduler)
}

func main() {
	assets := map[string]int64{
		"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d": time.Now().UnixNano(),
		// "0x4be3223f8708ca6b30d1e8b8926cf281ec83e770": time.Now().UnixNano(),
		// "0x8a90cab2b38dba80c64b7734e58ee1db38b8992e": time.Now().UnixNano(),
	}

	app := NewApp(assets)
	app.Schedule()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
