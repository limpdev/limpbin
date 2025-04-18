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
