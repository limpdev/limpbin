# `goco`

> icon picker written in golang!

```sh
[.ShellClassInfo]
IconResource=<absolute/path>,<#>
[ViewState]
Mode=
Vid=
FolderType=Generic
```

Here's a suggested approach and some Go code snippets to get you started:

**Core Logic: `desktop.ini` for Folders**

1.  **Understanding `desktop.ini`:**

    -   `desktop.ini` is a configuration file that allows you to customize how a folder is displayed in Windows Explorer.
    -   For custom icons, you'll use the `[.ShellClassInfo]` section and the `IconResource` key.
    -   `desktop.ini` must be located in the folder you want to customize.
    -   The folder itself must have the `System` and `ReadOnly` attributes set for `desktop.ini` to be effective.

2.  **Go Code Structure (Folder Customization):**

    ```go
    package main

    import (
        "fmt"
        "os"
        "path/filepath"
        "strings"
    )

    func main() {
        if len(os.Args) < 3 {
            fmt.Println("Usage: icon-picker <command> <target> [icon_path] [icon_index]")
            fmt.Println("Commands: set, remove")
            fmt.Println("Example (set folder icon): icon-picker set folder C:\\MyFolder C:\\icons\\myicon.ico 0")
            fmt.Println("Example (remove folder icon): icon-picker remove folder C:\\MyFolder")
            return
        }

        command := os.Args[1]
        targetPath := os.Args[2]

        switch command {
        case "set":
            if len(os.Args) < 5 {
                fmt.Println("Usage: icon-picker set folder <folder_path> <icon_path> <icon_index>")
                return
            }
            iconPath := os.Args[3]
            iconIndex := os.Args[4] // You'll need to validate and convert this to int later
            err := setFolderIcon(targetPath, iconPath, iconIndex)
            if err != nil {
                fmt.Println("Error setting folder icon:", err)
            } else {
                fmt.Println("Folder icon set successfully!")
            }

        case "remove":
            err := removeFolderIcon(targetPath)
            if err != nil {
                fmt.Println("Error removing folder icon:", err)
            } else {
                fmt.Println("Folder icon removed (default restored)!")
            }

        default:
            fmt.Println("Invalid command. Use 'set' or 'remove'.")
        }
    }

    func setFolderIcon(folderPath, iconPath, iconIndex string) error {
        desktopIniPath := filepath.Join(folderPath, "desktop.ini")
        iniContent := fmt.Sprintf(`[.ShellClassInfo]
        IconResource=%s,%s
        `, iconPath, iconIndex)

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

    // --- Placeholder functions for attributes (Platform-Specific!) ---

    func setFolderAttributes(folderPath string) error {
        // **Platform-Specific Implementation Needed (Windows)**
        // You'll need to use syscall to set file attributes on Windows.
        fmt.Println("Placeholder: Setting folder attributes (System, ReadOnly) for:", folderPath)
        return nil // Replace with actual Windows attribute setting logic
    }

    func clearFolderAttributes(folderPath string) error {
        // **Platform-Specific Implementation Needed (Windows)**
        // You'll need to use syscall to clear file attributes on Windows.
        fmt.Println("Placeholder: Clearing folder attributes (System, ReadOnly) for:", folderPath)
        return nil // Replace with actual Windows attribute clearing logic
    }
    ```

3.  **Platform-Specific Attribute Handling (Windows using `syscall`):**

    You'll need to use the `syscall` package in Go to interact with Windows API for setting file attributes. Here's how you'd implement `setFolderAttributes` and `clearFolderAttributes`:

    ```go
    import (
        // ... other imports ...
        "syscall"
    )

    const (
        FILE_ATTRIBUTE_READONLY = 0x00000001
        FILE_ATTRIBUTE_SYSTEM   = 0x00000004
    )

    func setFolderAttributes(folderPath string) error {
        pathPtr, err := syscall.UTF16PtrFromString(folderPath)
        if err != nil {
            return err
        }

        err = syscall.SetFileAttributes(pathPtr, FILE_ATTRIBUTE_SYSTEM|FILE_ATTRIBUTE_READONLY)
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
    ```

**Handling Standalone Files and Executables**

This is where it gets a bit more involved. You can't directly use `desktop.ini` for individual files and executables in the same way as folders. Here are the common approaches and how you might implement them with your command-line tool:

1.  **Shortcuts (.lnk files):** This is the most practical and user-friendly approach for setting custom icons for individual files and executables.

    -   **How it works:** You create a shortcut to the file/executable. Shortcuts can have their own custom icons set independently of the target file.
    -   **Implementation:**

        -   **Create/Modify Shortcuts:** You'll need a Go library or use `syscall` to create and modify Windows shortcut files (`.lnk`).  Libraries like `github.com/lxn/walk/walk` (though it's a UI library, it has shortcut creation functions) or `github.com/JamesStewy/go-win-shortcut` can help. Alternatively, you can use Windows Scripting Host (WSH) via COM automation from Go, but that's more complex.
        -   **Set Icon in Shortcut:**  When creating or modifying the shortcut, you'll specify the `IconLocation` property to point to your icon file and index.

    -   **Example Command:**

        ```bash
        icon-picker set file "C:\MyExecutable.exe" "C:\icons\myexeicon.ico" 0
        ```

    -   **Go Code Snippet (Conceptual - using a hypothetical shortcut library):**

        ```go
        // ... imports ...
        // Assuming you have a shortcut library like "shortcutlib"

        func setFileIcon(filePath, iconPath, iconIndex string) error {
            shortcutPath := filePath + ".lnk" // Or decide on a consistent shortcut naming scheme

            // Check if shortcut exists, if not, create one
            // ... (Shortcut library logic to create or load/modify shortcut) ...

            // Set the icon in the shortcut
            err := shortcutlib.SetShortcutIcon(shortcutPath, iconPath, iconIndex) // Hypothetical function
            if err != nil {
                return fmt.Errorf("failed to set shortcut icon: %w", err)
            }

            return nil
        }

        func removeFileIcon(filePath string) error {
            shortcutPath := filePath + ".lnk" // Assuming same naming scheme

            // Remove the shortcut icon (revert to default shortcut icon - depends on library)
            err := shortcutlib.RemoveShortcutIcon(shortcutPath) // Hypothetical function
            if err != nil {
                return fmt.Errorf("failed to remove shortcut icon: %w", err)
            }
            return nil
        }
        ```

2.  **Directly Modifying Executable Icons (More Complex):**

    -   **How it works:**  Executables (PE files) can embed icons within themselves. You _could_ potentially modify the icon resources inside an EXE.
    -   **Complexity:**  This is significantly more complex. You'd need to:
        -   Parse the PE file format.
        -   Locate the icon resources.
        -   Replace/modify the icon data.
        -   Rebuild the PE file (correctly updating checksums, etc.).
    -   **Go Libraries:** You'd likely need to use a library that understands PE file format and resource manipulation. This is a much more advanced approach and might be overkill for your initial tool.
    -   **Considerations:** Modifying executable files directly can be risky and might be flagged by security software. It's generally less recommended than using shortcuts for custom icons.

3.  **File Type Associations (System-Wide, Not Recommended for Individual Files):**

    -   You _could_ technically modify file type associations in the Windows Registry to change the default icon for _all_ files of a certain type (e.g., all `.exe` files).
    -   **Not Recommended:** This is a system-wide change and not suitable for customizing the icon of just _one_ specific file. It's also more invasive and requires registry manipulation, which can be error-prone and potentially harmful if done incorrectly.

**Recommended Approach: Focus on Folders and Shortcuts**

For your command-line tool, the most practical and user-friendly approach is to:

1.  **Implement `desktop.ini` customization for folders** (as shown in the initial Go code).
2.  **Use shortcuts to handle custom icons for standalone files and executables.**  This is the most common and expected way users customize icons for individual files in Windows.

**Next Steps and Improvements**

-   **Error Handling:** Add more robust error handling throughout your code.
-   **Input Validation:** Validate user inputs (e.g., check if paths exist, if icon index is a number).
-   **Icon Index Parsing:** Convert the `iconIndex` string argument to an integer using `strconv.Atoi`.
-   **Shortcut Library Integration:** Research and integrate a Go library for creating and modifying Windows shortcuts (e.g., `github.com/JamesStewy/go-win-shortcut`).
-   **User Feedback:** Provide clearer messages to the user about what actions are being taken.
-   **Help/Usage Information:** Improve the help output when the tool is run with `-h` or `--help` flags (using the `flag` package in Go can help with command-line argument parsing).
-   **Icon Browsing (Future Enhancement):** For a more complete "Choose Icon" replacement, you could add functionality to browse icon files (ICO, DLL, EXE) and display available icons in the command line (though this might get complex in a purely CLI environment).  Perhaps you could output a list of icons and their indices from a given file.

**Example Usage (with Shortcut approach for files):**

```bash
icon-picker set folder "C:\MyFolder" "C:\icons\folder_icon.ico" 0
icon-picker set file "C:\MyExecutable.exe" "C:\icons\exe_icon.ico" 0
icon-picker remove folder "C:\MyFolder"
icon-picker remove file "C:\MyExecutable.exe"
```

By starting with the folder `desktop.ini` logic and then adding shortcut support for files, you'll create a very useful command-line icon customization tool in Go! Remember to handle platform-specific aspects (like attributes and shortcut creation) carefully, especially when targeting Windows. Good luck! Let me know if you have more questions as you build this.

* * *

## `refinement #2`

> focusing on **folders** only

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	FILE_ATTRIBUTE_READONLY = 0x00000001
	FILE_ATTRIBUTE_SYSTEM   = 0x00000004
)

func setFolderIcon(folderPath, iconPath string, iconIndex int) error {
	desktopIniPath := filepath.Join(folderPath, "desktop.ini")
	iniContent := fmt.Sprintf(`[.ShellClassInfo]
IconResource=%s,%d
`, iconPath, iconIndex)

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

	err = syscall.SetFileAttributes(pathPtr, FILE_ATTRIBUTE_SYSTEM|FILE_ATTRIBUTE_READONLY)
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
```

**Key improvements in this version:**

1.  **Argument Parsing with `flag` package:**

    -   Uses the `flag` package for cleaner argument handling and help flag implementation.
    -   Handles `-h` and `--help` flags.
    -   Uses `flag.Args()` and `flag.Arg(index)` to access positional arguments, making the `main` function cleaner.

2.  **Input Validation:**

    -   **Argument Count:** Checks if the correct number of arguments (3) is provided.
    -   **Help Flag Check:** If `-h` or `--help` is used, it prints usage and exits.
    -   **Folder Path Validation:**
        -   Uses `os.Stat` to check if the target folder path exists.
        -   Checks if the path is actually a directory using `folderInfo.IsDir()`.
    -   **Icon File Path Validation:**
        -   Uses `os.Stat` to check if the icon file path exists.
    -   **Icon Index Validation:**
        -   Uses `strconv.Atoi` to convert the `iconIndexStr` to an integer.
        -   Checks for errors during conversion.
        -   Ensures the index is non-negative.

3.  **Clearer Error Messages:**

    -   Provides more informative error messages when validation fails, indicating the specific issue (e.g., folder not found, invalid index).

4.  **`printUsage` Function:**

    -   Encapsulates the usage information into a separate function for better organization and readability.
    -   Provides clear instructions, argument descriptions, and examples in the help output.

5.  **No "remove" command (as requested for this stage):**
    -   This version focuses solely on setting the folder icon as per the user's request to refine the "set" functionality first.  We can add "remove" functionality later if needed.

**To Compile and Run:**

1.  **Save:** Save the code as `icon-setter.go`.
2.  **Compile:** Open a terminal or command prompt in the directory where you saved the file and run:
    ```bash
    go build icon-setter.go
    ```
    This will create an executable file (e.g., `icon-setter.exe` on Windows, `icon-setter` on Linux/macOS).
3.  **Run:**
    ```bash
    icon-setter -h  # To see usage information
    icon-setter "C:\Path\To\Your\Folder" "C:\Path\To\Your\Icon.ico" 0
    icon-setter "C:\Another\Folder" "C:\Windows\System32\shell32.dll" 17
    ```
    Replace the paths and index with your actual folder, icon file, and icon index. Make sure to use double quotes around paths with spaces.

This refined version provides a solid command-line tool for setting folder icons with proper input validation and user guidance. Let me know if you'd like to add more features or refine it further!

* * *

## `iteration #3` from Claude

1.  Error Handling Enhancement:

-   Consider using a custom error type for domain-specific errors
-   Add more detailed error messages for filesystem operations
-   Implement proper cleanup in case of partial failures during icon setting

2.  Path Handling:

-   Add path sanitization to handle relative paths
-   Consider using filepath.Abs() to convert relative paths to absolute paths
-   Add checks for UNC paths which might behave differently

3.  Configuration and Constants:

```go
const (
    defaultPermissions = 0644
    desktopIniFile    = "desktop.ini"
    shellClassInfo    = "[.ShellClassInfo]"
)
```

4.  Safety Improvements:

-   Add validation for the icon file format (.ico, .dll, .icl)
-   Consider adding a backup/restore mechanism for existing desktop.ini files
-   Add validation for maximum icon index based on file type

5.  Feature Additions:

-   Add a dry-run flag to preview changes without applying them
-   Add a verbose mode for detailed logging
-   Add support for removing custom icons (the removeFolderIcon function exists but isn't exposed)
-   Add support for multiple folders in one operation

6.  Code Structure:

-   Split the code into separate packages (e.g., `icon`, `filesystem`)
-   Move Windows-specific code into a separate platform package
-   Add interfaces for better testability

7.  Testing:

-   Add unit tests for the core functionality
-   Add integration tests for filesystem operations
-   Add test fixtures for different icon file types

8.  Documentation:

-   Add godoc comments for exported functions
-   Include examples in the documentation
-   Document Windows version compatibility

Here's a sample implementation of some of these suggestions:

```go
package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
)

var (
    ErrInvalidFolder   = errors.New("invalid folder path")
    ErrInvalidIcon     = errors.New("invalid icon file")
    ErrInvalidIndex    = errors.New("invalid icon index")
    ErrOperationFailed = errors.New("operation failed")
)

type IconSetter struct {
    DryRun  bool
    Verbose bool
    Backup  bool
}

type IconFile struct {
    Path  string
    Index int
}

func NewIconSetter(dryRun, verbose, backup bool) *IconSetter {
    return &IconSetter{
        DryRun:  dryRun,
        Verbose: verbose,
        Backup:  backup,
    }
}

func (s *IconSetter) SetIcon(folderPath string, icon IconFile) error {
    // Validate and get absolute paths
    absFolder, err := filepath.Abs(folderPath)
    if err != nil {
        return fmt.Errorf("%w: %v", ErrInvalidFolder, err)
    }

    absIcon, err := filepath.Abs(icon.Path)
    if err != nil {
        return fmt.Errorf("%w: %v", ErrInvalidIcon, err)
    }

    // Validate folder
    if err := s.validateFolder(absFolder); err != nil {
        return err
    }

    // Validate icon
    if err := s.validateIcon(absIcon, icon.Index); err != nil {
        return err
    }

    if s.DryRun {
        s.log("Would set icon %s[%d] for folder %s", absIcon, icon.Index, absFolder)
        return nil
    }

    // Backup existing desktop.ini if needed
    if s.Backup {
        if err := s.backupDesktopIni(absFolder); err != nil {
            return fmt.Errorf("backup failed: %w", err)
        }
    }

    // Set the icon
    if err := setFolderIcon(absFolder, absIcon, icon.Index); err != nil {
        return fmt.Errorf("%w: %v", ErrOperationFailed, err)
    }

    s.log("Successfully set icon for %s", absFolder)
    return nil
}

func (s *IconSetter) validateFolder(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        return fmt.Errorf("%w: %v", ErrInvalidFolder, err)
    }
    if !info.IsDir() {
        return fmt.Errorf("%w: not a directory", ErrInvalidFolder)
    }
    return nil
}

func (s *IconSetter) validateIcon(path string, index int) error {
    if index < 0 {
        return fmt.Errorf("%w: negative index", ErrInvalidIndex)
    }

    ext := filepath.Ext(path)
    switch ext {
    case ".ico", ".dll", ".icl":
        // Valid extensions
    default:
        return fmt.Errorf("%w: unsupported file type %s", ErrInvalidIcon, ext)
    }

    return nil
}

func (s *IconSetter) backupDesktopIni(folderPath string) error {
    src := filepath.Join(folderPath, "desktop.ini")
    dst := filepath.Join(folderPath, "desktop.ini.bak")

    if _, err := os.Stat(src); err != nil {
        if os.IsNotExist(err) {
            return nil // Nothing to backup
        }
        return err
    }

    data, err := os.ReadFile(src)
    if err != nil {
        return err
    }

    return os.WriteFile(dst, data, 0644)
}

func (s *IconSetter) log(format string, args ...interface{}) {
    if s.Verbose {
        fmt.Printf(format+"\n", args...)
    }
}

func main() {
    var (
        dryRun  bool
        verbose bool
        backup  bool
        help    bool
    )

    flag.BoolVar(&dryRun, "dry-run", false, "Preview changes without applying them")
    flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
    flag.BoolVar(&backup, "backup", true, "Backup existing desktop.ini files")
    flag.BoolVar(&help, "help", false, "Show usage information")
    flag.BoolVar(&help, "h", false, "Show usage information (shorthand)")

    flag.Parse()

    if help {
        printUsage()
        return
    }

    args := flag.Args()
    if len(args) < 3 {
        fmt.Println("Error: Missing arguments.")
        printUsage()
        os.Exit(1)
    }

    iconIndex, err := strconv.Atoi(args[2])
    if err != nil {
        fmt.Printf("Error: Invalid icon index '%s'. Must be an integer.\n", args[2])
        os.Exit(1)
    }

    setter := NewIconSetter(dryRun, verbose, backup)
    icon := IconFile{
        Path:  args[1],
        Index: iconIndex,
    }

    if err := setter.SetIcon(args[0], icon); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

This improved version includes:

-   A more structured approach with a dedicated IconSetter type
-   Better error handling with custom error types
-   Support for dry-run and verbose modes
-   Automatic backup of existing desktop.ini files
-   Better path validation and handling
-   More extensive logging
-   Cleaner separation of concerns

The code is now more maintainable, testable, and robust while providing additional useful features for users. Would you like me to elaborate on any of these improvements or suggest additional enhancements?

---

## building a user interface...

Gio UI is a great choice for this project! Here's a suggested structure for your icon picker dialog:

```go
type IconPicker struct {
    // Window state
    theme   *material.Theme
    window  *app.Window

    // Icon management
    iconLibraries []string    // Paths to default libraries
    currentIcons  []IconData  // Current displayed icons
    selectedIcon  *IconData

    // UI state
    darkMode bool
    scale    float32
}

type IconData struct {
    path  string
    index int
    preview image.Image
}
```

Key features to implement:

1. Main Layout:
```go
// Grid of icons with customizable size
func (p *IconPicker) iconGrid(gtx layout.Context) layout.Dimensions {
    // Implement grid layout with larger preview sizes
    // Add hover effects for better UX
}
```

2. Library Management:
```go
func (p *IconPicker) loadLibraries() {
    // Load configured library paths from config file
    // Parse DLL/ICO files for icons
}

func (p *IconPicker) addLibrary(path string) {
    // Add new library to persistent config
}
```

3. Theme Support:
```go
func (p *IconPicker) toggleTheme() {
    p.darkMode = !p.darkMode
    // Update theme colors
}
```

4. Preview Panel:
```go
func (p *IconPicker) previewPanel(gtx layout.Context) layout.Dimensions {
    // Show larger preview of selected icon
    // Display icon details (path, index)
    // Add "Apply" button
}
```

Would you like me to develop any of these components further? We can start with whichever part you'd like to implement first.

Some additional features we could add:
- Icon search/filtering
- Quick preview on hover
- Recent icons section
- Icon size slider
- Export/import library configurations
