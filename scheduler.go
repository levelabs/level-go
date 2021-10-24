package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

func NewScheduler() *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)

	return s
}
