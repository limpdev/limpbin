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
