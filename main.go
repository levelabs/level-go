package main

import (
	"errors"
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

var (
	errCollectionManagerFailed = errors.New("collection manager failed to start")
)

type App struct {
	scheduler         *gocron.Scheduler
	collectionManager *c.CollectionManager

	cache *ristretto.Cache
	db    *badger.DB
}

func NewApp(assets map[string]int64) *App {
	scheduler := s.NewScheduler()

	cm, err := c.NewCollectionManager(assets)
	if err != nil {
		log.Fatal(errCollectionManagerFailed)
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
		scheduler:         scheduler,
		collectionManager: cm,
		cache:             cache,
		db:                db,
	}

	return &app
}

func (app *App) Schedule() {
	app.scheduler.Every(5).Seconds().Do(func() {
		asset, err := app.collectionManager.RunSequence()
		if err != nil {
			// do something
			fmt.Println(err)
			return
		}

		serialized, err := asset.MarshalJSON()
		if err != nil {
			// do something
		}

		address := asset.Address()
		err = app.db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(address), serialized)
			return err
		})

		if err != nil {
			fmt.Println("Db error", err)
		}
	})

	app.scheduler.Every(5).Seconds().Do(func() {
		err := app.db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"))
			if err != nil {
				// do something
			}

			var asset []byte
			err = item.Value(func(val []byte) error {
				asset = append([]byte{}, val...)
				return nil
			})

			fmt.Println("Asset in byte", string(asset))

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
