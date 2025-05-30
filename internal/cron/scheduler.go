package cron

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
)

func InitScheduleCron(job func()) {
	c := cron.New()
	err := c.AddFunc("0 0 * * *", job) // runs every day at 00:00
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	c.Start()
	log.Println("Cron scheduler started...")

	// Graceful shutdown handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down cron scheduler...")
	c.Stop()
}