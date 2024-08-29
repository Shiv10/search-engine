package utils

import (
	"fmt"
	"github.com/robfig/cron"
)

func StartCronJobs() {
	c := cron.New()
	c.AddFunc("0 * * * *", runEngine) // run job every hour
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs\n", cronCount)
}

func runEngine() {
	fmt.Println("Starting engine")
}