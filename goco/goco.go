package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	FILE_ATTRIBUTE_READONLY = 0x00000001
	FILE_ATTRIBUTE_SYSTEM   = 0x00000004
)

func setFolderIcon(folderPath, iconPath string, iconIndex int) error {
	desktopIniPath := filepath.Join(folderPath, "desktop.ini")
	iniContent := fmt.Sprintf(`[.ShellClassInfo]
IconResource=%s
iconIndex=%d
[ViewState]
FolderType=Generic`, iconPath, iconIndex)

	// Create or overwrite desktop.ini
	err := os.WriteFile(desktopIniPath, []byte(iniContent), 0644) // 0644 permissions (rw-r--r--)
	if err != nil {
		return fmt.Errorf("failed to write desktop.ini: %w", err)
	}

	// Set folder attributes (System and ReadOnly)
	err = setFolderAttributes(folderPath)
	if err != nil {
		// Clean up desktop.ini if setting attributes fails
		os.Remove(desktopIniPath)
		return fmt.Errorf("failed to set folder attributes: %w", err)
	}

	return nil
}

func removeFolderIcon(folderPath string) error {
	desktopIniPath := filepath.Join(folderPath, "desktop.ini")

	// Remove desktop.ini file
	err := os.Remove(desktopIniPath)
	if err != nil && !os.IsNotExist(err) { // Ignore "not exists" error
		return fmt.Errorf("failed to remove desktop.ini: %w", err)
	}

	// Remove System and ReadOnly attributes (optional - you might want to leave ReadOnly)
	err = clearFolderAttributes(folderPath) // Implement clearFolderAttributes
	if err != nil {
		return fmt.Errorf("failed to clear folder attributes: %w", err)
	}
	return nil
}

func setFolderAttributes(folderPath string) error {
	pathPtr, err := syscall.UTF16PtrFromString(folderPath)
	if err != nil {
		return err
	}

	err = syscall.SetFileAttributes(pathPtr, FILE_ATTRIBUTE_READONLY|FILE_ATTRIBUTE_SYSTEM)
	if err != nil {
		return err
	}
	return nil
}

func clearFolderAttributes(folderPath string) error {
	pathPtr, err := syscall.UTF16PtrFromString(folderPath)
	if err != nil {
		return err
	}

	// Remove System and ReadOnly attributes (by clearing them)
	err = syscall.SetFileAttributes(pathPtr, 0) // Setting to 0 clears all attributes
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "help", false, "Show usage information")
	flag.BoolVar(&helpFlag, "h", false, "Show usage information (shorthand)")

	flag.Parse()

	if helpFlag {
		printUsage()
		return
	}

	if len(flag.Args()) < 3 {
		fmt.Println("Error: Missing arguments.")
		printUsage()
		os.Exit(1)
	}

	targetFolderPath := flag.Arg(0)
	iconPath := flag.Arg(1)
	iconIndexStr := flag.Arg(2)

	// Validate target folder path
	folderInfo, err := os.Stat(targetFolderPath)
	if err != nil {
		fmt.Printf("Error: Target folder '%s' does not exist: %v\n", targetFolderPath, err)
		os.Exit(1)
	}
	if !folderInfo.IsDir() {
		fmt.Printf("Error: '%s' is not a folder.\n", targetFolderPath)
		os.Exit(1)
	}

	// Validate icon file path
	_, err = os.Stat(iconPath)
	if err != nil {
		fmt.Printf("Error: Icon file '%s' does not exist: %v\n", iconPath, err)
		os.Exit(1)
	}

	// Validate icon index
	iconIndex, err := strconv.Atoi(iconIndexStr)
	if err != nil {
		fmt.Printf("Error: Invalid icon index '%s'. Must be an integer.\n", iconIndexStr)
		os.Exit(1)
	}
	if iconIndex < 0 {
		fmt.Println("Error: Icon index must be a non-negative integer.")
		os.Exit(1)
	}

	err = setFolderIcon(targetFolderPath, iconPath, iconIndex)
	if err != nil {
		fmt.Println("Error setting folder icon:", err)
		os.Exit(1)
	}

	fmt.Println("Folder icon set successfully!")
}

func printUsage() {
	fmt.Println("Usage: icon-setter <folder_path> <icon_path> <icon_index>")
	fmt.Println()
	fmt.Println("Sets a custom icon for the specified folder using desktop.ini.")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <folder_path>   Path to the target folder.")
	fmt.Println("  <icon_path>     Path to the icon file (.ico, .dll, .icl).")
	fmt.Println("  <icon_index>    Index of the icon within the icon file. Use 0 for standalone .ico files.")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  icon-setter C:\\MyFolder C:\\icons\\myicon.ico 0")
	fmt.Println("  icon-setter C:\\AnotherFolder C:\\Windows\\System32\\shell32.dll 17")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help      Show this usage information.")
}
