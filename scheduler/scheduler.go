package scheduler

import (
	"github.com/go-co-op/gocron"
	"time"
)

func NewScheduler() *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)

	return s
}

func Start(s *gocron.Scheduler) {
	s.StartAsync()
}
