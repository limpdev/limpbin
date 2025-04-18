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
