package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
)

func startScheduler(cfg config.Config) {
	sch := gocron.NewScheduler(time.UTC)

	startEchoJob(cfg, sch)

	sch.StartAsync()
}

func startEchoJob(cfg config.Config, sch *gocron.Scheduler) {
	log.Info().Str("expression", cfg.Schedulers.EchoJob.Expression).Msg("Echo Job Attached")

	sch.Cron(cfg.Schedulers.EchoJob.Expression).Do(func() {
		log.Info().Msg("Echo Job Running...")
	})
}
