package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/night1010/everhealth/logger"
	"github.com/night1010/everhealth/repository"
	"github.com/robfig/cron"
)

func main() {
	logger.SetLogrusLogger()

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	c := cron.New()
	defer c.Stop()

	repo := repository.NewProductOrderRepository(db)

	background := context.Background()

	err = c.AddFunc("@hourly", func() {
		err := repo.CancelOrder(background)
		if err != nil {
			logger.Log.Error(err)
		}
	})
	if err != nil {
		logger.Log.Error(err)
	}

	err = c.AddFunc("@daily", func() {
		err := repo.ConfirmOrder(background)
		if err != nil {
			logger.Log.Error(err)
		}
	})
	if err != nil {
		logger.Log.Error(err)
	}
	go c.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
