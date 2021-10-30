package main

import (
	"errors"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/ristretto"
	"github.com/go-co-op/gocron"
	"github.com/levelabs/level-go/collection"
	s "github.com/levelabs/level-go/scheduler"
	"log"
	"net/http"
	"time"
)

const (
	collectionBAYC      = "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"
	colectionPartyDegen = "0x4be3223f8708ca6b30d1e8b8926cf281ec83e770"
)

var (
	errManagerFailed = errors.New("Manager failed to start")
)

type App struct {
	scheduler *gocron.Scheduler
	manager   *collection.Manager

	cache *ristretto.Cache
	db    *badger.DB
}

func NewApp(assets map[string]int64) *App {
	scheduler := s.NewScheduler()

	manager, err := collection.NewManager(assets)
	if err != nil {
		log.Fatal(errManagerFailed)
	}

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

	app := App{
		scheduler: scheduler,
		manager:   manager,
		cache:     cache,
		db:        db,
	}

	return &app
}

func (app *App) Schedule() {
	app.scheduler.Every(5).Seconds().Do(func() {
		asset, err := app.manager.RunSequence()
		if err != nil {
			// Handle Each Error Here!
			log.Print("[WARN]: Sequencing", err)
			return
		}

		log.Printf("[SUCCESS]: Asset completed sequencing %s", asset)

		serialized, err := asset.MarshalJSON()
		if err != nil {
			// Handle Each Error Here!
			return
		}

		address := asset.AddressBytes()
		err = app.db.Update(func(txn *badger.Txn) error {
			err := txn.Set(address, serialized)
			return err
		})

		if err != nil {
			log.Print("[ERROR]: Issue with DB", err)
		}
	})

	app.scheduler.Every(5).Seconds().Do(func() {
		err := app.db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"))
			if err != nil {
				// do something
			}

			var bytes []byte
			err = item.Value(func(val []byte) error {
				bytes = append([]byte{}, val...)
				return nil
			})

			log.Printf("[FOUND]: Asset in DB %s", bytes.Trait)

			return nil
		})

		if err != nil {
			// do something
		}
	})

	s.Start(app.scheduler)
}

func main() {
	assets := map[string]int64{
		collectionBAYC: time.Now().UnixNano(),
		// colectionPartyDegen: time.Now().UnixNano(),
		// "0x8a90cab2b38dba80c64b7734e58ee1db38b8992e": time.Now().UnixNano(),
	}

	app := NewApp(assets)
	app.Schedule()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
