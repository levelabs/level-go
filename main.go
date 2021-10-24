package main

import (
	"fmt"
	s "github.com/levelabs/level-go/scheduler"
	"log"
	"net/http"
)

func main() {
	scheduler := s.NewScheduler()

	scheduler.Every(5).Seconds().Do(func() {
		fmt.Println("Scheduler starting...")
	})

	s.Start(scheduler)

	connect()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
