package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Plist structures to parse Apple-style Info.plist manifests
type Plist struct {
	Dict PlistDict `xml:"dict"`
}

type PlistDict struct {
	Keys    []string `xml:"key"`
	Strings []string `xml:"string"`
}

type AppManifest struct {
	BundleID    string
	DisplayName string
	Executable  string
	Description string
	Path        string
}

func getBaseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "MicroDeck")
}

func ensureDirectories() {
	base := getBaseDir()
	dirs := []string{
		filepath.Join(base, "Apps"),
		filepath.Join(base, "Data"),
		filepath.Join(base, "Library"),
	}
	for _, d := range dirs {
		_ = os.MkdirAll(d, 0755)
	}
}

func parseManifest(plistPath string) (AppManifest, error) {
	data, err := os.ReadFile(plistPath)
	if err != nil {
		return AppManifest{}, err
	}

	var p Plist
	if err := xml.Unmarshal(data, &p); err != nil {
		return AppManifest{}, err
	}

	manifest := AppManifest{Path: filepath.Dir(plistPath)}
	for i := 0; i < len(p.Dict.Keys) && i < len(p.Dict.Strings); i++ {
		switch p.Dict.Keys[i] {
		case "CFBundleIdentifier":
			manifest.BundleID = p.Dict.Strings[i]
		case "CFBundleDisplayName":
			manifest.DisplayName = p.Dict.Strings[i]
		case "CFBundleExecutable":
			manifest.Executable = p.Dict.Strings[i]
		case "MicroDeckDescription":
			manifest.Description = p.Dict.Strings[i]
		}
	}
	return manifest, nil
}

func installApp(sourceDir string) error {
	ensureDirectories()
	plistPath := filepath.Join(sourceDir, "Info.plist")
	manifest, err := parseManifest(plistPath)
	if err != nil {
		return fmt.Errorf("invalid Info.plist: %w", err)
	}

	destDir := filepath.Join(getBaseDir(), "Apps", manifest.DisplayName+".microapp")
	_ = os.MkdirAll(destDir, 0755)

	// Copy Info.plist
	if err := copyFile(plistPath, filepath.Join(destDir, "Info.plist")); err != nil {
		return err
	}

	// Copy Executable
	execName := manifest.Executable
	if runtime.GOOS == "windows" && !strings.HasSuffix(execName, ".exe") {
		execName += ".exe"
	}

	srcExec := filepath.Join(sourceDir, execName)
	dstExec := filepath.Join(destDir, execName)
	if err := copyFile(srcExec, dstExec); err != nil {
		return fmt.Errorf("executable not found (%s): %w", srcExec, err)
	}
	_ = os.Chmod(dstExec, 0755)

	fmt.Printf("[OK] Installed '%s' (%s) to %s\n", manifest.DisplayName, manifest.BundleID, destDir)
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func listApps() []AppManifest {
	ensureDirectories()
	appsDir := filepath.Join(getBaseDir(), "Apps")
	entries, _ := os.ReadDir(appsDir)

	var apps []AppManifest
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".microapp") {
			plist := filepath.Join(appsDir, entry.Name(), "Info.plist")
			if m, err := parseManifest(plist); err == nil {
				apps = append(apps, m)
			}
		}
	}
	return apps
}

func launchInteractive() {
	reader := bufio.NewReader(os.Stdin)
	for {
		apps := listApps()
		fmt.Println("\n==========================================")
		fmt.Printf("  MicroDeck Launcher (%s/%s)\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("==========================================")

		if len(apps) == 0 {
			fmt.Println("No micro-apps installed. Run './start install <app>' first.")
			return
		}

		for i, app := range apps {
			fmt.Printf("  [%d] %s - %s\n", i+1, app.DisplayName, app.Description)
		}
		fmt.Println("  [0] Exit")
		fmt.Print("\nSelect an app to launch: ")

		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 0 || choice > len(apps) {
			fmt.Println("Invalid selection.")
			continue
		}
		if choice == 0 {
			break
		}

		selected := apps[choice-1]
		execName := selected.Executable
		if runtime.GOOS == "windows" && !strings.HasSuffix(execName, ".exe") {
			execName += ".exe"
		}

		execPath := filepath.Join(selected.Path, execName)
		dataDir := filepath.Join(getBaseDir(), "Data", selected.BundleID)
		_ = os.MkdirAll(dataDir, 0755)

		fmt.Printf("\n--- Launching %s ---\n", selected.DisplayName)
		cmd := exec.Command(execPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = append(os.Environ(), "MICRODECK_DATA_DIR="+dataDir)
		_ = cmd.Run()
		fmt.Printf("--- %s Exited ---\n", selected.DisplayName)
	}
}

func main() {
	if len(os.Args) < 2 {
		launchInteractive()
		return
	}

	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Usage: microdeck install <path_to_app_dir>")
			os.Exit(1)
		}
		if err := installApp(os.Args[2]); err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			os.Exit(1)
		}
	case "list":
		for _, app := range listApps() {
			fmt.Printf("- %s (%s): %s\n", app.DisplayName, app.BundleID, app.Description)
		}
	default:
		launchInteractive()
	}
}
