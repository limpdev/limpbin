```go
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"goplunk/internal/app"
	"goplunk/internal/config"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start the application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the application
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
	defer application.Stop()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down...")
}
```