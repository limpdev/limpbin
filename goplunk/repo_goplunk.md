This file is a merged representation of the entire codebase, combined into a single document by Repomix.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded

## Additional Info

# Directory Structure
```
go.mod
internal/animation/animation.go
internal/app.go
internal/app/app.go
internal/config/config.go
internal/input/input.go
internal/renderer/renderer.go
main.go
```

# Files

## File: go.mod
```
module goplunk

go 1.24.2

require golang.org/x/sys v0.32.0
```

## File: internal/animation/animation.go
```go
package animation

import (
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"
)

func New(cfg *config.Config) (*App, error) {
	// Seed random number generator (important for particle animation)
	rand.Seed(time.Now().UnixNano())
	return nil, nil // Placeholder return to fix the function
}

// Animation is the interface for all animations
type Animation interface {
	Update(float32) bool // returns true if animation is still active
	GetPosition() (float32, float32)
}

// AnimationManager manages multiple animations
type AnimationManager struct {
	animations map[int]Animation
	nextID     int
	mutex      sync.RWMutex
	maxCount   int
}

// NewAnimationManager creates a new animation manager
func NewAnimationManager(maxCount int) *AnimationManager {
	return &AnimationManager{
		animations: make(map[int]Animation),
		nextID:     1,
		maxCount:   maxCount,
	}
}

// AddAnimation adds a new animation
func (am *AnimationManager) AddAnimation(anim Animation) int {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Check if we're at the max count
	if len(am.animations) >= am.maxCount {
		// Remove oldest animation
		oldestID := -1
		for id := range am.animations {
			if oldestID == -1 || id < oldestID {
				oldestID = id
			}
		}

		if oldestID != -1 {
			delete(am.animations, oldestID)
		}
	}

	id := am.nextID
	am.animations[id] = anim
	am.nextID++

	return id
}

// RemoveAnimation removes an animation by ID
func (am *AnimationManager) RemoveAnimation(id int) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	delete(am.animations, id)
}

// Update updates all animations
func (am *AnimationManager) Update(deltaTime float32) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Update all animations and remove finished ones
	for id, anim := range am.animations {
		if !anim.Update(deltaTime) {
			delete(am.animations, id)
		}
	}
}

// GetActiveAnimations returns all active animations
func (am *AnimationManager) GetActiveAnimations() []Animation {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	animations := make([]Animation, 0, len(am.animations))
	for _, anim := range am.animations {
		animations = append(animations, anim)
	}

	return animations
}

// RippleAnimation represents a ripple effect animation
type RippleAnimation struct {
	X, Y          float32
	CurrentRadius float32
	MaxRadius     float32
	Duration      float32
	ElapsedTime   float32
	Color         color.RGBA
	StartTime     time.Time
}

// NewRippleAnimation creates a new ripple animation
func NewRippleAnimation(x, y float32, maxRadius float32, duration float32, color color.RGBA) *RippleAnimation {
	return &RippleAnimation{
		X:             x,
		Y:             y,
		CurrentRadius: 0,
		MaxRadius:     maxRadius,
		Duration:      duration,
		ElapsedTime:   0,
		Color:         color,
		StartTime:     time.Now(),
	}
}

// Update updates the ripple animation
func (ra *RippleAnimation) Update(deltaTime float32) bool {
	ra.ElapsedTime += deltaTime

	// Calculate progress (0 to 1)
	progress := ra.ElapsedTime / ra.Duration
	if progress >= 1.0 {
		return false // Animation completed
	}

	// Update radius using easing function
	ra.CurrentRadius = ra.MaxRadius * easeOutCubic(progress)

	// Update alpha color based on progress (fade out)
	ra.Color.A = uint8(255 * (1.0 - progress))

	return true
}

// GetPosition returns the position of the animation
func (ra *RippleAnimation) GetPosition() (float32, float32) {
	return ra.X, ra.Y
}

// ParticleAnimation represents a particle effect animation
type ParticleAnimation struct {
	X, Y        float32
	Particles   []Particle
	Duration    float32
	ElapsedTime float32
	Color       color.RGBA
}

// Particle represents a single particle in a particle animation
type Particle struct {
	X, Y        float32
	VelocityX   float32
	VelocityY   float32
	Size        float32
	Rotation    float32
	RotationVel float32
	Alpha       float32
}

// NewParticleAnimation creates a new particle animation
func NewParticleAnimation(x, y float32, particleCount int, duration float32, color color.RGBA) *ParticleAnimation {
	pa := &ParticleAnimation{
		X:           x,
		Y:           y,
		Duration:    duration,
		ElapsedTime: 0,
		Color:       color,
		Particles:   make([]Particle, particleCount),
	}

	// Initialize particles with random velocities and properties
	for i := 0; i < particleCount; i++ {
		angle := float32(2 * math.Pi * float64(i) / float64(particleCount))
		speed := float32(40.0 + 20.0*math.Rand()) // Random speed between 40-60

		pa.Particles[i] = Particle{
			X:           x,
			Y:           y,
			VelocityX:   speed * float32(math.Cos(float64(angle))),
			VelocityY:   speed * float32(math.Sin(float64(angle))),
			Size:        float32(3.0 + 5.0*math.Rand()), // Random size between 3-8
			Rotation:    float32(2 * math.Pi * math.Rand()),
			RotationVel: float32(math.Pi * math.Rand() * 2), // Random rotation velocity
			Alpha:       1.0,
		}
	}

	return pa
}

// Update updates the particle animation
func (pa *ParticleAnimation) Update(deltaTime float32) bool {
	pa.ElapsedTime += deltaTime

	// Calculate progress (0 to 1)
	progress := pa.ElapsedTime / pa.Duration
	if progress >= 1.0 {
		return false // Animation completed
	}

	// Update particles
	for i := range pa.Particles {
		p := &pa.Particles[i]

		// Update position
		p.X += p.VelocityX * deltaTime
		p.Y += p.VelocityY * deltaTime

		// Apply gravity
		p.VelocityY += 50.0 * deltaTime

		// Apply friction
		p.VelocityX *= 0.98
		p.VelocityY *= 0.98

		// Update rotation
		p.Rotation += p.RotationVel * deltaTime

		// Fade out
		p.Alpha = 1.0 - progress
	}

	return true
}

// GetPosition returns the position of the animation
func (pa *ParticleAnimation) GetPosition() (float32, float32) {
	return pa.X, pa.Y
}

// SplashAnimation represents a splash effect animation
type SplashAnimation struct {
	X, Y        float32
	Rings       []float32 // Ring radiuses
	Duration    float32
	ElapsedTime float32
	Color       color.RGBA
}

// NewSplashAnimation creates a new splash animation
func NewSplashAnimation(x, y float32, maxRadius float32, duration float32, color color.RGBA) *SplashAnimation {
	return &SplashAnimation{
		X:           x,
		Y:           y,
		Rings:       []float32{0, 0, 0}, // Three rings with different speeds
		Duration:    duration,
		ElapsedTime: 0,
		Color:       color,
	}
}

// Update updates the splash animation
func (sa *SplashAnimation) Update(deltaTime float32) bool {
	sa.ElapsedTime += deltaTime

	// Calculate progress (0 to 1)
	progress := sa.ElapsedTime / sa.Duration
	if progress >= 1.0 {
		return false // Animation completed
	}

	// Update ring radiuses at different speeds
	sa.Rings[0] = 100.0 * easeOutQuad(progress)
	sa.Rings[1] = 70.0 * easeOutQuad(float32(math.Max(0, float64(progress)-0.1)))
	sa.Rings[2] = 40.0 * easeOutQuad(float32(math.Max(0, float64(progress)-0.2)))

	// Update alpha color based on progress (fade out)
	sa.Color.A = uint8(255 * (1.0 - progress))

	return true
}

// GetPosition returns the position of the animation
func (sa *SplashAnimation) GetPosition() (float32, float32) {
	return sa.X, sa.Y
}

// Easing functions
func easeOutQuad(t float32) float32 {
	return t * (2 - t)
}

func easeOutCubic(t float32) float32 {
	t = t - 1
	return t*t*t + 1
}
```

## File: internal/app.go
```go
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
```

## File: internal/app/app.go
```go
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
```

## File: internal/config/config.go
```go
package config

import (
	"encoding/json"
	"errors"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
)

// AnimationStyle defines the type of animation to use
type AnimationStyle int

const (
	Ripple AnimationStyle = iota
	Particles
	Splash
)

// Config holds the application configuration
type Config struct {
	// Animation properties
	Color            color.RGBA
	AnimationDuration float32
	MaxConcurrent     int
	Style             AnimationStyle
	
	// Rendering properties
	FrameRate        int
	ParticleCount    int
	MaxRadius        float32
	
	// Behavior properties
	EnableStartup     bool
	ShowTrayIcon      bool
	CaptureMiddleClick bool
	CaptureRightClick  bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Color:            color.RGBA{59, 130, 246, 255}, // Blue
		AnimationDuration: 0.8,
		MaxConcurrent:     10,
		Style:             Ripple,
		FrameRate:         60,
		ParticleCount:     15,
		MaxRadius:         100.0,
		EnableStartup:     true,
		ShowTrayIcon:      true,
		CaptureMiddleClick: false,
		CaptureRightClick:  false,
	}
}

// hexToRGBA converts a hex color string to color.RGBA
func hexToRGBA(hex string) (color.RGBA, error) {
	if hex[0] == '#' {
		hex = hex[1:]
	}
	
	if len(hex) != 6 {
		return color.RGBA{}, errors.New("invalid hex color format")
	}
	
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

// Load loads configuration from the specified file, falling back to default
func Load(configFile string) (*Config, error) {
	cfg := DefaultConfig()
	
	// If no config file is specified, return default
	if configFile == "" {
		// Try to find config in standard locations
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configFile = filepath.Join(homeDir, ".config", "goplunk", "config.json")
			if _, err := os.Stat(configFile); err != nil {
				// Config doesn't exist, return default
				return cfg, nil
			}
		} else {
			// Can't determine home directory, use default
			return cfg, nil
		}
	}
	
	// Open and parse config file
	file, err := os.Open(configFile)
	if err != nil {
		return cfg, err
	}
	defer file.Close()
	
	var jsonCfg struct {
		Color            string  `json:"color"`
		AnimationDuration float32  `json:"animation_duration"`
		MaxConcurrent     int     `json:"max_concurrent"`
		Style             string  `json:"style"`
		FrameRate         int     `json:"frame_rate"`
		ParticleCount     int     `json:"particle_count"`
		MaxRadius         float32 `json:"max_radius"`
		EnableStartup     bool    `json:"enable_startup"`
		ShowTrayIcon      bool    `json:"show_tray_icon"`
		CaptureMiddleClick bool    `json:"capture_middle_click"`
		CaptureRightClick  bool    `json:"capture_right_click"`
	}
	
	if err := json.NewDecoder(file).Decode(&jsonCfg); err != nil {
		return cfg, err
	}
	
	// Apply values from JSON config
	if jsonCfg.Color != "" {
		color, err := hexToRGBA(jsonCfg.Color)
		if err == nil {
			cfg.Color = color
		}
	}
	
	if jsonCfg.AnimationDuration > 0 {
		cfg.AnimationDuration = jsonCfg.AnimationDuration
	}
	
	if jsonCfg.MaxConcurrent > 0 {
		cfg.MaxConcurrent = jsonCfg.MaxConcurrent
	}
	
	if jsonCfg.FrameRate > 0 {
		cfg.FrameRate = jsonCfg.FrameRate
	}
	
	if jsonCfg.ParticleCount > 0 {
		cfg.ParticleCount = jsonCfg.ParticleCount
	}
	
	if jsonCfg.MaxRadius > 0 {
		cfg.MaxRadius = jsonCfg.MaxRadius
	}
	
	// Set animation style
	switch jsonCfg.Style {
	case "ripple":
		cfg.Style = Ripple
	case "particles":
		cfg.Style = Particles
	case "splash":
		cfg.Style = Splash
	}
	
	cfg.EnableStartup = jsonCfg.EnableStartup
	cfg.ShowTrayIcon = jsonCfg.ShowTrayIcon
	cfg.CaptureMiddleClick = jsonCfg.CaptureMiddleClick
	cfg.CaptureRightClick = jsonCfg.CaptureRightClick
	
	return cfg, nil
}

// Save saves the configuration to the specified file
func (c *Config) Save(configFile string) error {
	// Ensure directory exists
	dir := filepath.Dir(configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Create or truncate file
	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Convert RGBA color to hex
	hexColor := "#" + 
		strconv.FormatUint(uint64(c.Color.R), 16) +
		strconv.FormatUint(uint64(c.Color.G), 16) +
		strconv.FormatUint(uint64(c.Color.B), 16)
	
	// Create JSON struct
	jsonCfg := struct {
		Color            string  `json:"color"`
		AnimationDuration float32  `json:"animation_duration"`
		MaxConcurrent     int     `json:"max_concurrent"`
		Style             string  `json:"style"`
		FrameRate         int     `json:"frame_rate"`
		ParticleCount     int     `json:"particle_count"`
		MaxRadius         float32 `json:"max_radius"`
		EnableStartup     bool    `json:"enable_startup"`
		ShowTrayIcon      bool    `json:"show_tray_icon"`
		CaptureMiddleClick bool    `json:"capture_middle_click"`
		CaptureRightClick  bool    `json:"capture_right_click"`
	}{
		Color:            hexColor,
		AnimationDuration: c.AnimationDuration,
		MaxConcurrent:     c.MaxConcurrent,
		FrameRate:         c.FrameRate,
		ParticleCount:     c.ParticleCount,
		MaxRadius:         c.MaxRadius,
		EnableStartup:     c.EnableStartup,
		ShowTrayIcon:      c.ShowTrayIcon,
		CaptureMiddleClick: c.CaptureMiddleClick,
		CaptureRightClick:  c.CaptureRightClick,
	}
	
	// Set style string
	switch c.Style {
	case Ripple:
		jsonCfg.Style = "ripple"
	case Particles:
		jsonCfg.Style = "particles"
	case Splash:
		jsonCfg.Style = "splash"
	}
	
	// Encode to JSON with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonCfg)
}
```

## File: internal/input/input.go
```go
package input

import (
	"fmt"
	"log"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Event types
type EventType int

const (
	LeftClick EventType = iota
	RightClick
	MiddleClick
)

// MouseEvent represents a mouse click event
type MouseEvent struct {
	Type EventType
	X    int32
	Y    int32
}

// MouseHook manages the Windows mouse hooks
type MouseHook struct {
	hookHandle  windows.Handle
	subscribers []chan MouseEvent
	mutex       sync.RWMutex
	running     bool
}

// Windows constants
const (
	WH_MOUSE_LL    = 14
	WM_LBUTTONDOWN = 0x0201
	WM_RBUTTONDOWN = 0x0204
	WM_MBUTTONDOWN = 0x0207
)

// MSLLHOOKSTRUCT represents the Windows low-level mouse hook struct
type MSLLHOOKSTRUCT struct {
	Pt          windows.Point
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

var (
	user32                = windows.NewLazyDLL("user32.dll")
	setWindowsHookEx      = user32.NewProc("SetWindowsHookExW")
	callNextHookEx        = user32.NewProc("CallNextHookEx")
	unhookWindowsHookEx   = user32.NewProc("UnhookWindowsHookEx")
	getModuleHandle       = windows.NewLazyDLL("kernel32.dll").NewProc("GetModuleHandleW")
	mouseHookCallback     uintptr
	globalMouseHookInstance *MouseHook
)

// hookCallback is the Windows hook callback function
func hookCallback(code int, wParam uintptr, lParam uintptr) uintptr {
	if code >= 0 {
		var eventType EventType
		
		switch wParam {
		case WM_LBUTTONDOWN:
			eventType = LeftClick
		case WM_RBUTTONDOWN:
			eventType = RightClick
		case WM_MBUTTONDOWN:
			eventType = MiddleClick
		default:
			goto callNext
		}
		
		// Extract mouse coordinates from MSLLHOOKSTRUCT
		mouseStruct := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		
		// Notify subscribers
		if globalMouseHookInstance != nil {
			globalMouseHookInstance.notify(MouseEvent{
				Type: eventType,
				X:    mouseStruct.Pt.X,
				Y:    mouseStruct.Pt.Y,
			})
		}
	}
	
callNext:
	ret, _, _ := callNextHookEx.Call(0, uintptr(code), wParam, lParam)
	return ret
}

// NewMouseHook creates a new MouseHook instance
func NewMouseHook() (*MouseHook, error) {
	if globalMouseHookInstance != nil {
		return nil, fmt.Errorf("mouse hook already installed")
	}
	
	hook := &MouseHook{
		subscribers: make([]chan MouseEvent, 0),
		running:     false,
	}
	
	// Set global instance for callback access
	globalMouseHookInstance = hook
	
	return hook, nil
}

// Subscribe adds a subscriber channel
func (h *MouseHook) Subscribe(ch chan MouseEvent) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.subscribers = append(h.subscribers, ch)
}

// Unsubscribe removes a subscriber channel
func (h *MouseHook) Unsubscribe(ch chan MouseEvent) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	for i, subscriber := range h.subscribers {
		if subscriber == ch {
			h.subscribers = append(h.subscribers[:i], h.subscribers[i+1:]...)
			break
		}
	}
}

// notify sends an event to all subscribers
func (h *MouseHook) notify(event MouseEvent) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	for _, ch := range h.subscribers {
		select {
		case ch <- event:
		default:
			// Channel is full, skip
		}
	}
}

// Start installs the mouse hook
func (h *MouseHook) Start() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if h.running {
		return nil
	}
	
	// Create a callback function
	mouseHookCallback = windows.NewCallback(hookCallback)
	
	// Get module handle for the current process (NULL)
	moduleHandle, _, _ := getModuleHandle.Call(0)
	
	// Set Windows hook
	hookHandle, _, err := setWindowsHookEx.Call(
		WH_MOUSE_LL,
		mouseHookCallback,
		moduleHandle,
		0,
	)
	
	if hookHandle == 0 {
		return fmt.Errorf("failed to set mouse hook: %v", err)
	}
	
	h.hookHandle = windows.Handle(hookHandle)
	h.running = true
	log.Println("Mouse hook installed successfully")
	
	return nil
}

// Stop uninstalls the mouse hook
func (h *MouseHook) Stop() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if !h.running {
		return nil
	}
	
	ret, _, err := unhookWindowsHookEx.Call(uintptr(h.hookHandle))
	if ret == 0 {
		return fmt.Errorf("failed to unhook mouse hook: %v", err)
	}
	
	h.running = false
	h.hookHandle = 0
	globalMouseHookInstance = nil
	log.Println("Mouse hook uninstalled successfully")
	
	return nil
}
```

## File: internal/renderer/renderer.go
```go
// internal/renderer.go
package renderer

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"goplunk/internal/animation" // Assuming animation is moved to internal
	"goplunk/internal/config"    // Assuming config is moved to internal

	"golang.org/x/sys/windows"
)

// Direct2D libraries and functions (assuming these are defined globally or loaded lazily)
var (
	user32   = windows.NewLazyDLL("user32.dll")
	kernel32 = windows.NewLazyDLL("kernel32.dll")
	d2d1     = windows.NewLazyDLL("d2d1.dll")
	// Function pointers
	getClassInfoEx             = user32.NewProc("GetClassInfoExW")
	registerClassEx            = user32.NewProc("RegisterClassExW")
	unregisterClass            = user32.NewProc("UnregisterClassW")
	createWindowEx             = user32.NewProc("CreateWindowExW")
	destroyWindow              = user32.NewProc("DestroyWindow")
	defWindowProc              = user32.NewProc("DefWindowProcW")
	getSystemMetrics           = user32.NewProc("GetSystemMetrics")
	setLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	peekMessage                = user32.NewProc("PeekMessageW")
	translateMessage           = user32.NewProc("TranslateMessage")
	dispatchMessage            = user32.NewProc("DispatchMessageW")
	postQuitMessage            = user32.NewProc("PostQuitMessage")
	getModuleHandle            = kernel32.NewProc("GetModuleHandleW")
	d2d1CreateFactory          = d2d1.NewProc("D2D1CreateFactory")
)

// Direct2D interfaces (simplified representation using uintptr)
type (
	IUnknown              uintptr
	ID2D1Factory          uintptr
	ID2D1HwndRenderTarget uintptr
	ID2D1SolidColorBrush  uintptr
	ID2D1StrokeStyle      uintptr
)

// Direct2D constants (add missing ones)
const (
	// Window styles
	WS_EX_LAYERED     = 0x80000
	WS_EX_TRANSPARENT = 0x20
	WS_EX_NOACTIVATE  = 0x08000000
	WS_EX_TOOLWINDOW  = 0x00000080
	WS_POPUP          = 0x80000000
	WS_VISIBLE        = 0x10000000
	// Layered window attributes
	LWA_ALPHA = 0x2
	// System metrics
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
	// Window messages
	WM_DESTROY = 0x0002
	WM_PAINT   = 0x000F
	// PeekMessage options
	PM_REMOVE = 0x0001
	// D2D constants
	D2D1_FACTORY_TYPE_SINGLE_THREADED = 0
	D2D1_FEATURE_LEVEL_DEFAULT        = 0
	D2D1_PRESENT_OPTIONS_NONE         = 0x0
	D2D1_RENDER_TARGET_TYPE_DEFAULT   = 0
	D2D1_RENDER_TARGET_USAGE_NONE     = 0
	D2D1_PIXEL_FORMAT_UNKNOWN         = 0
	D2D1_ALPHA_MODE_PREMULTIPLIED     = 1
	D2D1_CAP_STYLE_ROUND              = 1
	D2D1_LINE_JOIN_ROUND              = 1
	D2D1_DASH_STYLE_SOLID             = 0
)

// --- Helper Functions & Structs ---

// releaseCOM releases a COM object safely
func releaseCOM(obj IUnknown) {
	if obj != 0 {
		syscall.SyscallN(uintptr(obj)+2*unsafe.Sizeof(uintptr(0)), uintptr(obj)) // Call Release method (IUnknown index 2)
	}
}

// RECT structure for Windows API
type RECT struct {
	Left, Top, Right, Bottom int32
}

// POINT structure for Windows API
type POINT struct {
	X, Y int32
}

// MSG structure for Windows API
type MSG struct {
	Hwnd    windows.HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// WNDCLASSEX structure for Windows API
type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     windows.HINSTANCE
	HIcon         windows.HICON
	HCursor       windows.HCURSOR
	HbrBackground windows.HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.HICON
}

// D2D1_SIZE_U structure
type D2D1_SIZE_U struct {
	Width, Height uint32
}

// D2D1_PIXEL_FORMAT structure
type D2D1_PIXEL_FORMAT struct {
	Format    uint32 // DXGI_FORMAT enum value (use D2D1_PIXEL_FORMAT_UNKNOWN for default)
	AlphaMode uint32 // D2D1_ALPHA_MODE enum value
}

// D2D1_RENDER_TARGET_PROPERTIES structure
type D2D1_RENDER_TARGET_PROPERTIES struct {
	Type        uint32            // D2D1_RENDER_TARGET_TYPE
	PixelFormat D2D1_PIXEL_FORMAT // D2D1_PIXEL_FORMAT
	DpiX        float32
	DpiY        float32
	Usage       uint32 // D2D1_RENDER_TARGET_USAGE
	MinLevel    uint32 // D2D1_FEATURE_LEVEL
}

// D2D1_HWND_RENDER_TARGET_PROPERTIES structure
type D2D1_HWND_RENDER_TARGET_PROPERTIES struct {
	Hwnd           windows.HWND
	PixelSize      D2D1_SIZE_U
	PresentOptions uint32 // D2D1_PRESENT_OPTIONS
}

// D2D1_COLOR_F structure
type D2D1_COLOR_F struct {
	R, G, B, A float32
}

// D2D1_POINT_2F structure
type D2D1_POINT_2F struct {
	X, Y float32
}

// D2D1_ELLIPSE structure
type D2D1_ELLIPSE struct {
	Point   D2D1_POINT_2F
	RadiusX float32
	RadiusY float32
}

// D2D1_RECT_F structure
type D2D1_RECT_F struct {
	Left, Top, Right, Bottom float32
}

// D2D1_STROKE_STYLE_PROPERTIES structure
type D2D1_STROKE_STYLE_PROPERTIES struct {
	StartCap   uint32 // D2D1_CAP_STYLE
	EndCap     uint32 // D2D1_CAP_STYLE
	DashCap    uint32 // D2D1_CAP_STYLE
	LineJoin   uint32 // D2D1_LINE_JOIN
	MiterLimit float32
	DashStyle  uint32 // D2D1_DASH_STYLE
	DashOffset float32
}

// IID_ID2D1Factory GUID for Direct2D factory
var IID_ID2D1Factory = windows.GUID{
	Data1: 0x06152247, Data2: 0x6f50, Data3: 0x465a, Data4: [8]byte{0x92, 0x45, 0x11, 0x8b, 0xfd, 0x3b, 0x60, 0x07},
}

// --- Renderer Implementation ---

// Renderer handles Direct2D rendering
type Renderer struct {
	hwnd             windows.HWND
	factory          ID2D1Factory
	renderTarget     ID2D1HwndRenderTarget
	brushes          map[color.RGBA]ID2D1SolidColorBrush
	strokeStyle      ID2D1StrokeStyle
	animationManager *animation.AnimationManager
	config           *config.Config
	width, height    int32
	mutex            sync.Mutex
	running          bool
	stopEvent        chan struct{} // Used to signal render loop to stop
	windowClassName  *uint16
	hInstance        windows.HINSTANCE
}

// NewRenderer creates a new Direct2D renderer
func NewRenderer(animManager *animation.AnimationManager, cfg *config.Config) (*Renderer, error) {
	runtime.LockOSThread() // Direct2D and windowing often require thread affinity

	r := &Renderer{
		brushes:          make(map[color.RGBA]ID2D1SolidColorBrush),
		animationManager: animManager,
		config:           cfg,
		stopEvent:        make(chan struct{}),
	}

	// 1. Get Module Handle
	hInst, _, err := getModuleHandle.Call(0)
	if hInst == 0 {
		return nil, fmt.Errorf("failed to get module handle: %w", err)
	}
	r.hInstance = windows.HINSTANCE(hInst)

	// 2. Create D2D Factory
	var pFactory ID2D1Factory
	ret, _, err := d2d1CreateFactory.Call(
		uintptr(D2D1_FACTORY_TYPE_SINGLE_THREADED),
		uintptr(unsafe.Pointer(&IID_ID2D1Factory)),
		0, // D2D1_FACTORY_OPTIONS not used
		uintptr(unsafe.Pointer(&pFactory)),
	)
	if int32(ret) < 0 { // Check HRESULT for failure
		return nil, fmt.Errorf("failed to create D2D factory (HRESULT: 0x%X): %w", uint32(ret), err)
	}
	r.factory = pFactory
	log.Println("Direct2D Factory created")

	// 3. Get Screen Dimensions
	w, _, _ := getSystemMetrics.Call(uintptr(SM_CXSCREEN))
	h, _, _ := getSystemMetrics.Call(uintptr(SM_CYSCREEN))
	r.width = int32(w)
	r.height = int32(h)

	// 4. Register Window Class
	className := "GoPlunkOverlayWindow"
	r.windowClassName, _ = windows.UTF16PtrFromString(className)

	var wc WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = windows.NewCallback(wndProc)
	wc.HInstance = r.hInstance
	wc.LpszClassName = r.windowClassName
	wc.Style = 0x0008 // CS_OWNDC for Direct2D performance

	ret, _, err = registerClassEx.Call(uintptr(unsafe.Pointer(&wc)))
	if ret == 0 {
		// Check if already registered (might happen on restart)
		ret, _, _ := getClassInfoEx.Call(uintptr(r.hInstance), uintptr(unsafe.Pointer(r.windowClassName)), uintptr(unsafe.Pointer(&wc)))
		if ret == 0 {
			return nil, fmt.Errorf("failed to register window class: %w", err)
		}
		log.Println("Window class already registered, reusing.")
	} else {
		log.Println("Window class registered")
	}

	// 5. Create Window
	hwnd, _, err := createWindowEx.Call(
		uintptr(WS_EX_LAYERED|WS_EX_TRANSPARENT|WS_EX_NOACTIVATE|WS_EX_TOOLWINDOW),
		uintptr(unsafe.Pointer(r.windowClassName)),
		0, // No window title
		uintptr(WS_POPUP|WS_VISIBLE),
		0, 0, // Position (x, y) - covering full screen
		uintptr(r.width), uintptr(r.height), // Size
		0, // No parent window
		0, // No menu
		uintptr(r.hInstance),
		0, // No extra parameter
	)
	if hwnd == 0 {
		releaseCOM(IUnknown(r.factory))                                                        // Clean up factory
		unregisterClass.Call(uintptr(unsafe.Pointer(r.windowClassName)), uintptr(r.hInstance)) // Clean up class
		return nil, fmt.Errorf("failed to create window: %w", err)
	}
	r.hwnd = windows.HWND(hwnd)
	log.Printf("Overlay window created (HWND: %X)", hwnd)

	// Set initial transparency (fully transparent) - might not be strictly necessary
	// setLayeredWindowAttributes.Call(hwnd, 0, 255, LWA_ALPHA)

	// 6. Create D2D Render Target
	renderTargetProps := D2D1_RENDER_TARGET_PROPERTIES{
		Type:        D2D1_RENDER_TARGET_TYPE_DEFAULT,
		PixelFormat: D2D1_PIXEL_FORMAT{Format: D2D1_PIXEL_FORMAT_UNKNOWN, AlphaMode: D2D1_ALPHA_MODE_PREMULTIPLIED},
		DpiX:        0, // Use system DPI
		DpiY:        0,
		Usage:       D2D1_RENDER_TARGET_USAGE_NONE,
		MinLevel:    D2D1_FEATURE_LEVEL_DEFAULT,
	}
	hwndRenderTargetProps := D2D1_HWND_RENDER_TARGET_PROPERTIES{
		Hwnd:           r.hwnd,
		PixelSize:      D2D1_SIZE_U{Width: uint32(r.width), Height: uint32(r.height)},
		PresentOptions: D2D1_PRESENT_OPTIONS_NONE,
	}

	var pRenderTarget ID2D1HwndRenderTarget
	ret, _, err = syscall.Syscall6(r.factory+8*unsafe.Sizeof(uintptr(0)), // VTable index 8 for CreateHwndRenderTarget
		4, // Number of arguments
		uintptr(r.factory),
		uintptr(unsafe.Pointer(&renderTargetProps)),
		uintptr(unsafe.Pointer(&hwndRenderTargetProps)),
		uintptr(unsafe.Pointer(&pRenderTarget)),
		0, 0)
	if int32(ret) < 0 {
		destroyWindow.Call(hwnd)
		releaseCOM(IUnknown(r.factory))
		unregisterClass.Call(uintptr(unsafe.Pointer(r.windowClassName)), uintptr(r.hInstance))
		return nil, fmt.Errorf("failed to create D2D render target (HRESULT: 0x%X): %w", uint32(ret), err)
	}
	r.renderTarget = pRenderTarget
	log.Println("Direct2D Render Target created")

	// 7. Create Stroke Style
	strokeStyleProps := D2D1_STROKE_STYLE_PROPERTIES{
		StartCap:   D2D1_CAP_STYLE_ROUND,
		EndCap:     D2D1_CAP_STYLE_ROUND,
		DashCap:    D2D1_CAP_STYLE_ROUND, // Relevant for dashed lines
		LineJoin:   D2D1_LINE_JOIN_ROUND,
		MiterLimit: 10.0,
		DashStyle:  D2D1_DASH_STYLE_SOLID,
		DashOffset: 0.0,
	}
	var pStrokeStyle ID2D1StrokeStyle
	ret, _, err = syscall.Syscall6(r.factory+18*unsafe.Sizeof(uintptr(0)), // VTable index 18 for CreateStrokeStyle
		4, // Number of arguments
		uintptr(r.factory),
		uintptr(unsafe.Pointer(&strokeStyleProps)),
		0, // Dashes array (nullptr for solid)
		0, // Dashes count
		uintptr(unsafe.Pointer(&pStrokeStyle)), 0)
	if int32(ret) < 0 {
		releaseCOM(IUnknown(r.renderTarget))
		destroyWindow.Call(hwnd)
		releaseCOM(IUnknown(r.factory))
		unregisterClass.Call(uintptr(unsafe.Pointer(r.windowClassName)), uintptr(r.hInstance))
		return nil, fmt.Errorf("failed to create D2D stroke style (HRESULT: 0x%X): %w", uint32(ret), err)
	}
	r.strokeStyle = pStrokeStyle
	log.Println("Direct2D Stroke Style created")

	return r, nil
}

// Start begins the rendering loop
func (r *Renderer) Start() {
	r.mutex.Lock()
	if r.running {
		r.mutex.Unlock()
		return
	}
	r.running = true
	r.mutex.Unlock()

	go r.renderLoop()
	go r.messageLoop() // Need a message loop for the window
	log.Println("Renderer started")
}

// Stop cleans up resources and stops rendering
func (r *Renderer) Stop() {
	r.mutex.Lock()
	if !r.running {
		r.mutex.Unlock()
		return
	}
	r.running = false
	close(r.stopEvent) // Signal loops to stop
	r.mutex.Unlock()

	// Post a quit message to break the message loop if it's blocking
	postQuitMessage.Call(0)

	// Give loops a moment to exit (optional, adjust as needed)
	time.Sleep(100 * time.Millisecond)

	runtime.LockOSThread() // Ensure cleanup happens on the correct thread
	defer runtime.UnlockOSThread()

	log.Println("Stopping renderer and cleaning up resources...")
	r.mutex.Lock() // Re-acquire lock for cleanup
	defer r.mutex.Unlock()

	// Release brushes first
	for _, brush := range r.brushes {
		releaseCOM(IUnknown(brush))
	}
	r.brushes = make(map[color.RGBA]ID2D1SolidColorBrush) // Clear map

	releaseCOM(IUnknown(r.strokeStyle))
	r.strokeStyle = 0
	releaseCOM(IUnknown(r.renderTarget))
	r.renderTarget = 0
	releaseCOM(IUnknown(r.factory))
	r.factory = 0

	if r.hwnd != 0 {
		destroyWindow.Call(uintptr(r.hwnd))
		r.hwnd = 0
	}

	if r.windowClassName != nil && r.hInstance != 0 {
		unregisterClass.Call(uintptr(unsafe.Pointer(r.windowClassName)), uintptr(r.hInstance))
		r.windowClassName = nil // Prevent double unregister
	}

	log.Println("Renderer stopped and resources released.")
}

// messageLoop handles window messages
func (r *Renderer) messageLoop() {
	runtime.LockOSThread() // Ensure message loop runs on the window's thread
	defer runtime.UnlockOSThread()

	var msg MSG
	for {
		// Use PeekMessage to avoid blocking indefinitely, allowing stopEvent check
		ret, _, _ := peekMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, PM_REMOVE)
		if ret != 0 { // If message retrieved
			if msg.Message == 0x0012 { // WM_QUIT
				log.Println("Received WM_QUIT, exiting message loop.")
				break // Exit loop if WM_QUIT received
			}
			translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}

		// Check if stop is requested
		select {
		case <-r.stopEvent:
			log.Println("Stop event received, exiting message loop.")
			return
		default:
			// Prevent busy-waiting, yield CPU slightly
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// renderLoop continuously draws frames
func (r *Renderer) renderLoop() {
	runtime.LockOSThread() // Ensure rendering happens on the correct thread
	defer runtime.UnlockOSThread()

	targetFrameTime := time.Second / time.Duration(r.config.FrameRate)
	ticker := time.NewTicker(targetFrameTime)
	defer ticker.Stop()

	lastTime := time.Now()

	for {
		select {
		case <-r.stopEvent:
			log.Println("Stop event received, exiting render loop.")
			return
		case <-ticker.C:
			now := time.Now()
			deltaTime := float32(now.Sub(lastTime).Seconds())
			lastTime = now

			// Update animations (outside draw lock)
			r.animationManager.Update(deltaTime)

			// Draw frame
			r.drawFrame()
		}
	}
}

// drawFrame renders a single frame with all active animations
func (r *Renderer) drawFrame() {
	r.mutex.Lock() // Lock during drawing operations
	defer r.mutex.Unlock()

	if !r.running || r.renderTarget == 0 {
		return // Don't draw if not running or target is invalid
	}

	// --- Begin Draw ---
	syscall.Syscall(r.renderTarget+3*unsafe.Sizeof(uintptr(0)), 1, uintptr(r.renderTarget), 0, 0) // VTable index 3: BeginDraw

	// --- Clear Background ---
	clearColor := D2D1_COLOR_F{R: 0, G: 0, B: 0, A: 0}             // Transparent black
	syscall.Syscall(r.renderTarget+5*unsafe.Sizeof(uintptr(0)), 2, // VTable index 5: Clear
		uintptr(r.renderTarget),
		uintptr(unsafe.Pointer(&clearColor)),
		0)

	// --- Draw Animations ---
	activeAnimations := r.animationManager.GetActiveAnimations()
	for _, anim := range activeAnimations {
		switch a := anim.(type) {
		case *animation.RippleAnimation:
			r.drawRipple(a)
		case *animation.ParticleAnimation:
			r.drawParticles(a)
		case *animation.SplashAnimation:
			r.drawSplash(a)
		}
	}

	// --- End Draw ---
	ret, _, _ := syscall.Syscall(r.renderTarget+4*unsafe.Sizeof(uintptr(0)), 1, uintptr(r.renderTarget), 0, 0) // VTable index 4: EndDraw
	if int32(ret) < 0 {
		// Handle potential errors, e.g., D2DERR_RECREATE_TARGET
		log.Printf("EndDraw failed (HRESULT: 0x%X), potential device loss", uint32(ret))
		// Basic handling: log error. More robust: recreate target.
		// For simplicity here, we just log. In a real app, you might need to call Stop() and restart/recreate renderer.
	}
}

// getBrush retrieves or creates a brush for the given color
func (r *Renderer) getBrush(clr color.RGBA) (ID2D1SolidColorBrush, error) {
	// Assumes r.mutex is already held by drawFrame

	if brush, ok := r.brushes[clr]; ok && brush != 0 {
		return brush, nil
	}

	// Create new brush
	d2dColor := D2D1_COLOR_F{
		R: float32(clr.R) / 255.0,
		G: float32(clr.G) / 255.0,
		B: float32(clr.B) / 255.0,
		A: float32(clr.A) / 255.0, // Direct2D expects non-premultiplied here
	}

	var pBrush ID2D1SolidColorBrush
	ret, _, err := syscall.Syscall(r.renderTarget+15*unsafe.Sizeof(uintptr(0)), // VTable index 15: CreateSolidColorBrush
		3, // Number of arguments
		uintptr(r.renderTarget),
		uintptr(unsafe.Pointer(&d2dColor)),
		0, // Brush properties (optional, can be nil)
		uintptr(unsafe.Pointer(&pBrush)), 0, 0)

	if int32(ret) < 0 {
		return 0, fmt.Errorf("failed to create solid color brush (HRESULT: 0x%X): %w", uint32(ret), err)
	}

	r.brushes[clr] = pBrush // Store the new brush
	return pBrush, nil
}

// --- Drawing Functions ---

func (r *Renderer) drawRipple(a *animation.RippleAnimation) {
	if a.Color.A == 0 {
		return
	} // Skip invisible

	brush, err := r.getBrush(a.Color)
	if err != nil || brush == 0 {
		log.Printf("Error getting brush for ripple: %v", err)
		return
	}

	ellipse := D2D1_ELLIPSE{
		Point:   D2D1_POINT_2F{X: a.X, Y: a.Y},
		RadiusX: a.CurrentRadius,
		RadiusY: a.CurrentRadius,
	}
	strokeWidth := float32(2.0) // Example stroke width

	syscall.Syscall6(r.renderTarget+9*unsafe.Sizeof(uintptr(0)), // VTable index 9: DrawEllipse
		5, // Number of arguments
		uintptr(r.renderTarget),
		uintptr(unsafe.Pointer(&ellipse)),
		uintptr(brush),
		uintptr(math.Float32bits(strokeWidth)), // Pass float as uintptr
		uintptr(r.strokeStyle), 0)
}

func (r *Renderer) drawParticles(a *animation.ParticleAnimation) {
	// Reuse a single base color brush and modify its opacity for each particle
	baseColor := a.Color
	baseColor.A = 255 // Use full alpha for the base brush lookup/creation
	brush, err := r.getBrush(baseColor)
	if err != nil || brush == 0 {
		log.Printf("Error getting brush for particles: %v", err)
		return
	}

	originalOpacity, _, _ := syscall.Syscall(brush+4*unsafe.Sizeof(uintptr(0)), 1, uintptr(brush), 0, 0) // VTable index 4: GetOpacity

	for i := range a.Particles {
		p := &a.Particles[i]
		if p.Alpha <= 0 {
			continue
		}

		particleAlpha := p.Alpha * (float32(a.Color.A) / 255.0) // Combine particle alpha and animation alpha
		if particleAlpha <= 0 {
			continue
		}

		// Set brush opacity for this particle
		syscall.Syscall(brush+5*unsafe.Sizeof(uintptr(0)), 2, // VTable index 5: SetOpacity
			uintptr(brush),
			uintptr(math.Float32bits(particleAlpha)), // Pass float as uintptr
			0)

		// Draw particle (simple filled ellipse)
		ellipse := D2D1_ELLIPSE{
			Point:   D2D1_POINT_2F{X: p.X, Y: p.Y},
			RadiusX: p.Size / 2.0,
			RadiusY: p.Size / 2.0,
		}
		syscall.Syscall(r.renderTarget+10*unsafe.Sizeof(uintptr(0)), // VTable index 10: FillEllipse
			3, // Number of arguments
			uintptr(r.renderTarget),
			uintptr(unsafe.Pointer(&ellipse)),
			uintptr(brush), 0, 0, 0)
	}

	// Restore original brush opacity
	syscall.Syscall(brush+5*unsafe.Sizeof(uintptr(0)), 2, uintptr(brush), originalOpacity, 0)
}

func (r *Renderer) drawSplash(a *animation.SplashAnimation) {
	if a.Color.A == 0 {
		return
	} // Skip invisible

	brush, err := r.getBrush(a.Color)
	if err != nil || brush == 0 {
		log.Printf("Error getting brush for splash: %v", err)
		return
	}

	strokeWidth := float32(1.5) // Example stroke width

	for _, radius := range a.Rings {
		if radius <= 0 {
			continue
		}
		ellipse := D2D1_ELLIPSE{
			Point:   D2D1_POINT_2F{X: a.X, Y: a.Y},
			RadiusX: radius,
			RadiusY: radius,
		}

		syscall.Syscall6(r.renderTarget+9*unsafe.Sizeof(uintptr(0)), // VTable index 9: DrawEllipse
			5, // Number of arguments
			uintptr(r.renderTarget),
			uintptr(unsafe.Pointer(&ellipse)),
			uintptr(brush),
			uintptr(math.Float32bits(strokeWidth)),
			uintptr(r.strokeStyle), 0)
	}
}

// --- Window Procedure ---

// wndProc handles messages for the overlay window
func wndProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_DESTROY:
		log.Println("WM_DESTROY received, posting quit message.")
		postQuitMessage.Call(0) // Signal termination
		return 0
	case WM_PAINT:
		// We handle drawing ourselves in the render loop, so just validate
		syscall.SyscallN(user32.NewProc("ValidateRect").Addr(), uintptr(hwnd), 0) // Avoid default painting
		return 0
	default:
		// Let Windows handle other messages
		ret, _, _ := defWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return ret
	}
}
```

## File: main.go
```go
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
```
