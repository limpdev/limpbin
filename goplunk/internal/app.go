// internal/app/app.go
package app

import (
	"fmt"
	"log"
	"math/rand" // Import math/rand
	"time"      // Import time

	"goplunk/internal/animation"
	"goplunk/internal/config"
	"goplunk/internal/input"
	"goplunk/internal/renderer"
)

// App encapsulates the application state
type App struct {
	config           *config.Config
	animationManager *animation.AnimationManager
	mouseHook        *input.MouseHook
	renderer         *renderer.Renderer
	eventChan        chan input.MouseEvent
	stopChan         chan struct{} // For signaling goroutines to stop
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Seed random number generator (important for particle animation)
	rand.Seed(time.Now().UnixNano())

	animManager := animation.NewAnimationManager(cfg.MaxConcurrent)

	hook, err := input.NewMouseHook()
	if err != nil {
		// Don't treat "already installed" as fatal if we might be restarting
		log.Printf("Warning creating mouse hook: %v. Application might still work if hook is already running.", err)
		// return nil, fmt.Errorf("failed to create mouse hook: %w", err)
	}

	rend, err := renderer.NewRenderer(animManager, cfg)
	if err != nil {
		// Attempt cleanup if renderer failed after hook potentially succeeded
		if hook != nil {
			_ = hook.Stop() // Ignore error on cleanup
		}
		return nil, fmt.Errorf("failed to initialize renderer: %w", err)
	}

	return &App{
		config:           cfg,
		animationManager: animManager,
		mouseHook:        hook, // May be nil if NewMouseHook failed but we continued
		renderer:         rend,
		eventChan:        make(chan input.MouseEvent, 32), // Buffered channel
		stopChan:         make(chan struct{}),
	}, nil
}

// Start initializes and runs the application components
func (a *App) Start() error {
	log.Println("Starting application...")

	// Subscribe to mouse events *before* starting the hook
	if a.mouseHook != nil {
		a.mouseHook.Subscribe(a.eventChan)
		if err := a.mouseHook.Start(); err != nil {
			a.renderer.Stop() // Stop renderer if hook fails
			return fmt.Errorf("failed to start mouse hook: %w", err)
		}
	} else {
		log.Println("Mouse hook not available, click detection disabled.")
	}

	a.renderer.Start()

	// Start event processing goroutine
	go a.processEvents()

	log.Println("Application started successfully.")
	return nil
}

// Stop cleans up application resources
func (a *App) Stop() {
	log.Println("Stopping application...")
	close(a.stopChan) // Signal event processor to stop

	if a.mouseHook != nil {
		// Unsubscribe first to prevent sending to closed channel
		// (though buffered channel makes this less critical)
		a.mouseHook.Unsubscribe(a.eventChan)
		if err := a.mouseHook.Stop(); err != nil {
			log.Printf("Error stopping mouse hook: %v", err)
		}
	}

	if a.renderer != nil {
		a.renderer.Stop()
	}

	// Close event channel after ensuring producer (hook) and consumer (processEvents) are stopped
	// Check if already closed by processEvents potentially
	select {
	case _, ok := <-a.eventChan:
		if ok {
			close(a.eventChan)
		}
	default:
		// Channel was empty or already closed
		close(a.eventChan)
	}

	log.Println("Application stopped.")
}

// processEvents listens for mouse clicks and triggers animations
func (a *App) processEvents() {
	log.Println("Event processor started.")
	defer log.Println("Event processor stopped.")

	for {
		select {
		case event, ok := <-a.eventChan:
			if !ok {
				log.Println("Event channel closed.")
				return // Channel closed, exit goroutine
			}

			// Check config if we should handle this click type
			shouldHandle := false
			switch event.Type {
			case input.LeftClick:
				shouldHandle = true // Always handle left click for now
			case input.MiddleClick:
				shouldHandle = a.config.CaptureMiddleClick
			case input.RightClick:
				shouldHandle = a.config.CaptureRightClick
			}

			if shouldHandle {
				// log.Printf("Click detected: Type=%d, Pos=(%d, %d)", event.Type, event.X, event.Y)
				var anim animation.Animation
				x, y := float32(event.X), float32(event.Y)

				switch a.config.Style {
				case config.Ripple:
					anim = animation.NewRippleAnimation(x, y, a.config.MaxRadius, a.config.AnimationDuration, a.config.Color)
				case config.Particles:
					anim = animation.NewParticleAnimation(x, y, a.config.ParticleCount, a.config.AnimationDuration, a.config.Color)
				case config.Splash:
					anim = animation.NewSplashAnimation(x, y, a.config.MaxRadius, a.config.AnimationDuration, a.config.Color)
				default:
					log.Printf("Unknown animation style: %d", a.config.Style)
					continue // Skip if style is unknown
				}

				a.animationManager.AddAnimation(anim)
			}

		case <-a.stopChan:
			log.Println("Stop signal received by event processor.")
			return // Stop signal received, exit goroutine
		}
	}
}
