// main.go (or goplunk/main.go)
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath" // Import filepath
	"runtime"
	"syscall"

	"goplunk/internal/app"    // Import the app package
	"goplunk/internal/config" // Import the local config package
)

func main() {
	// Ensure the main goroutine runs on the main OS thread.
	// This is often required for GUI/windowing operations on some platforms,
	// although renderer.go tries to handle thread affinity itself.
	// It's generally good practice for the main entry point.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread() // Unlock shouldn't strictly be necessary on exit, but good practice

	// --- Configuration Loading ---
	// Default config path
	homeDir, err := os.UserHomeDir()
	defaultConfigPath := ""
	if err == nil {
		defaultConfigPath = filepath.Join(homeDir, ".config", "goplunk", "config.json")
	}

	// Parse command line flags
	configFile := flag.String("config", defaultConfigPath, "Path to config file (defaults to ~/.config/goplunk/config.json)")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		// If loading failed specifically because the default file doesn't exist, it's not fatal.
		// We proceed with the default config object returned by Load in that case.
		if !os.IsNotExist(err) || *configFile != defaultConfigPath {
			log.Printf("Warning: Failed to load configuration from %s: %v. Using defaults.", *configFile, err)
		} else {
			log.Printf("No config file found at %s. Using default settings.", defaultConfigPath)
		}
		// If Load returned defaults despite error (e.g., file not found), cfg might still be valid.
		// Let's double-check cfg is not nil. If Load strictly returns nil on error, use DefaultConfig().
		if cfg == nil {
			log.Println("Falling back to hardcoded default config.")
			cfg = config.DefaultConfig()
		}
	} else {
		log.Printf("Configuration loaded from %s", *configFile)
	}

	// Optionally save the loaded/default config back if it didn't exist
	// This helps users discover the config file format and location.
	if *configFile != "" && os.IsNotExist(err) { // Only if Load failed due to non-existence
		log.Printf("Creating default config file at %s", *configFile)
		if saveErr := cfg.Save(*configFile); saveErr != nil {
			log.Printf("Warning: Failed to save default config file: %v", saveErr)
		}
	}

	// --- Application Setup ---
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the application
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// --- Signal Handling & Shutdown ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Application running. Press Ctrl+C to exit.")

	// Wait for termination signal
	<-sigChan
	log.Println("Shutdown signal received...")

	// Stop the application gracefully
	application.Stop()

	log.Println("GoPlunk exited.")
}
