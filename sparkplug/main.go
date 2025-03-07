package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed static/*
var content embed.FS

func main() {
	// Command 1: hugo serve
	hugoCmd := exec.Command("hugo", "serve", "--disableFastRender")
	hugoCmd.Stdout = os.Stdout // Pipe output to console
	hugoCmd.Stderr = os.Stderr // Pipe errors to console

	fmt.Println("Starting Hugo server...")
	if err := hugoCmd.Start(); err != nil {
		fmt.Println("Error starting Hugo server:", err)
		return
	}

	// Command 2: tailscale serve 1313
	tailscaleCmd := exec.Command("tailscale", "serve", "1313")
	tailscaleCmd.Stdout = os.Stdout // Pipe output to console
	tailscaleCmd.Stderr = os.Stderr // Pipe errors to console

	fmt.Println("Starting Tailscale serve on port 1313...")
	if err := tailscaleCmd.Start(); err != nil {
		fmt.Println("Error starting Tailscale serve:", err)
		// Clean up Hugo process if Tailscale fails to start - not strictly necessary for this example but good practice
		if err := hugoCmd.Process.Kill(); err != nil {
			fmt.Println("Error killing Hugo process:", err)
		}
		return
	}

	// Wait for both commands to finish (though Hugo serve will run indefinitely until stopped)
	if err := hugoCmd.Wait(); err != nil {
		fmt.Println("Hugo server stopped with error:", err)
	}
	if err := tailscaleCmd.Wait(); err != nil {
		fmt.Println("Tailscale serve stopped with error:", err)
	}

	fmt.Println("Automation process finished.")
}
