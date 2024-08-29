package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"search-engine/db"
	"search-engine/routes"
	"search-engine/utils"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func main() {
	env := godotenv.Load()

	if env != nil {
		panic("cannot find environment variabled")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = ":4000"
	} else {
		port = ":" + port
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())
	db.InitDB()
	routes.SetRoutes(app)
	utils.StartCronJobs()
	// Start server and listen for shutdown
	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // block the main thread
	app.Shutdown()
	fmt.Println("Shutting down server")
}
