package main

import (
	// "fmt"
	// "github.com/go-co-op/gocron"
	"github.com/levelabs/level-go/collection"
	// s "github.com/levelabs/level-go/scheduler"
	// "log"
	// "net/http"
)

// const (
// 	collectionBAYC      = "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"
// 	colectionPartyDegen = "0x4be3223f8708ca6b30d1e8b8926cf281ec83e770"
// )
//
// type App struct {
// 	scheduler  *gocron.Scheduler
// 	collection *Collections
// }
//
// func NewApp() *App {
// 	scheduler := s.NewScheduler()
// 	app := App{scheduler: scheduler}
// 	return &app
// }

func main() {
	collection.CollectionPriorityTest()

	// NewApp()
	//
	// c := collection.NewAsset(colectionPartyDegen)
	// c.SetBaseURI()
	// c.QueryAttributes()
	//
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal(err)
	// }
}
