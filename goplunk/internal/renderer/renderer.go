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
