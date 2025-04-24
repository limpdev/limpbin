## `modules\animation-go.go`

```go
package animation

import (
	"image/color"
	"math"
	"sync"
	"time"
)

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
	X, Y           float32
	CurrentRadius  float32
	MaxRadius      float32
	Duration       float32
	ElapsedTime    float32
	Color          color.RGBA
	StartTime      time.Time
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
	X, Y         float32
	Particles    []Particle
	Duration     float32
	ElapsedTime  float32
	Color        color.RGBA
}

// Particle represents a single particle in a particle animation
type Particle struct {
	X, Y         float32
	VelocityX    float32
	VelocityY    float32
	Size         float32
	Rotation     float32
	RotationVel  float32
	Alpha        float32
}

// NewParticleAnimation creates a new particle animation
func NewParticleAnimation(x, y float32, particleCount int, duration float32, color color.RGBA) *ParticleAnimation {
	pa := &ParticleAnimation{
		X:          x,
		Y:          y,
		Duration:   duration,
		ElapsedTime: 0,
		Color:      color,
		Particles:  make([]Particle, particleCount),
	}
	
	// Initialize particles with random velocities and properties
	for i := 0; i < particleCount; i++ {
		angle := float32(2 * math.Pi * float64(i) / float64(particleCount))
		speed := float32(40.0 + 20.0*math.Rand()) // Random speed between 40-60
		
		pa.Particles[i] = Particle{
			X:          x,
			Y:          y,
			VelocityX:  speed * float32(math.Cos(float64(angle))),
			VelocityY:  speed * float32(math.Sin(float64(angle))),
			Size:       float32(3.0 + 5.0*math.Rand()), // Random size between 3-8
			Rotation:   float32(2 * math.Pi * math.Rand()),
			RotationVel: float32(math.Pi * math.Rand() * 2), // Random rotation velocity
			Alpha:      1.0,
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
	X, Y           float32
	Rings          []float32 // Ring radiuses
	Duration       float32
	ElapsedTime    float32
	Color          color.RGBA
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
	sa.Rings[1] = 70.0 * easeOutQuad(math.Max(0, progress-0.1))
	sa.Rings[2] = 40.0 * easeOutQuad(math.Max(0, progress-0.2))
	
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

## `modules\config-go.go`

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

## `modules\input-go.go`

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

## `modules\main-go.go`

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

## `modules\renderer.go`

```go
package renderer

import (
	"fmt"
	"image/color"
	"log"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
	"goplunk/internal/animation"
)

// Direct2D libraries
var (
	d2d1              = windows.NewLazyDLL("d2d1.dll")
	dwriteLib         = windows.NewLazyDLL("dwrite.dll")
	createFactory     = d2d1.NewProc("D2D1CreateFactory")
	releaseCOM        = windows.NewProc("Release")
)

// Direct2D interfaces
type (
	ID2D1Factory            uintptr
	ID2D1HwndRenderTarget   uintptr
	ID2D1SolidColorBrush    uintptr
	ID2D1StrokeStyle        uintptr
)

// Direct2D constants
const (
	D2D1_FACTORY_TYPE_SINGLE_THREADED    = 0
	D2D1_RENDER_TARGET_TYPE_DEFAULT      = 0
	D2D1_FEATURE_LEVEL_DEFAULT           = 0
	D2D1_HWND_RENDER_TARGET_PROPERTIES   = 0x2
	D2D1_STROKE_STYLE_PROPERTIES         = 0x8
	CLSCTX_INPROC_SERVER                 = 0x1
	D2D1_CAP_STYLE_ROUND                 = 1
	D2D1_STROKE_STYLE_PROPERTIES_DASH_STYLE_SOLID = 0
	D2D1_ALPHA_MODE_PREMULTIPLIED        = 1
	WS_EX_LAYERED                        = 0x80000
	WS_EX_TRANSPARENT                    = 0x20
	WS_EX_NOACTIVATE                     = 0x08000000
	WS_EX_TOOLWINDOW                     = 0x00000080
	WS_POPUP                             = 0x80000000
	WS_VISIBLE                           = 0x10000000
	LWA_ALPHA                            = 0x2
)

// Renderer handles Direct2D rendering
type Renderer struct {
	hwnd                 windows.HWND
	factory              ID2D1Factory
	renderTarget         ID2D1HwndRenderTarget
	brushes              map[color.RGBA]ID2D1SolidColorBrush
	strokeStyle          ID2D1StrokeStyle
	animationManager     *animation.AnimationManager
	width, height        int32
	mutex                sync.Mutex
	running              bool
	frameRate            int
}

// RECT structure for Windows API
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// D2D1_SIZE_U structure
type D2D1_SIZE_U struct {
	Width  uint32
	Height uint32
}

// D2D1_PIXEL_FORMAT structure
type D2D1_PIXEL_FORMAT struct {
	Format    uint32
	AlphaMode uint32
}

// D2D1_HWND_RENDER_TARGET_PROPERTIES structure
type D2D1_HWND_RENDER_TARGET_PROPERTIES struct {
	Hwnd           windows.HWND
	PixelSize      D2D1_SIZE_U
	PresentOptions uint32
}

// D2D1_RENDER_TARGET_PROPERTIES structure
type D2D1_RENDER_TARGET_PROPERTIES struct {
	Type         uint32
	PixelFormat  D2D1_PIXEL_FORMAT
	DpiX         float32
	DpiY         float32
	Usage        uint32
	MinLevel     uint32
}

// D2D1_COLOR_F structure
type D2D1_COLOR_F struct {
	R float32
	G float32
	B float32
	A float32
}

// D2D1_ELLIPSE structure
type D2D1_ELLIPSE struct {
	X     float32
	Y     float32
	RadiusX float32
	RadiusY float32
}

// D2D1_STROKE_STYLE_PROPERTIES structure
type D2D1_STROKE_STYLE_PROPERTIES struct {
	StartCap     uint32
	EndCap       uint32
	DashCap      uint32
	LineJoin     uint32
	MiterLimit   float32
	DashStyle    uint32
	DashOffset   float32
}

// IID_ID2D1Factory GUID for Direct2D factory
var IID_ID2D1Factory = windows.GUID{
	Data1: 0x06152247,
	Data2: 0x6f50,
	Data3: 0x465a,
	Data4: [8]byte{0x92, 0x45, 0x11, 0x8b, 0xfd, 0x3b, 0x60, 0x07},
}

// NewRenderer creates a new Direct2D renderer
func NewRenderer(animManager *animation.AnimationManager, frameRate int) (*Renderer, error) {
	r := &Renderer{
		brushes:          make(map[color.RGBA
// HERE -> LEFT OFF TO FINISH LATER
// CONTINUE FROM HERE
```

