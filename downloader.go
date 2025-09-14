package ffmpeg_downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// IsSystemInstalled checks if FFmpeg is installed and accessible system-wide.
//
// # Returns:
//   - bool: true if FFmpeg is found and can be executed system-wide, otherwise false.
//
// # Example:
//
//	if IsSystemInstalled() {
//	    // FFmpeg is available system-wide
//	    cmd := exec.Command("ffmpeg", "-i", "input.mp4", "output.mp4")
//	} else {
//	    // Need to download or use statically installed version
//	    path := GetFFmpegPath("my-app")
//	}
func IsSystemInstalled() bool {
	_, err := exec.Command("ffmpeg", "-version").CombinedOutput()
	if err != nil {
		return false
	}

	return true
}

// IsStaticallyInstalled checks if FFmpeg is installed in the specified configuration directory.
//
// # Parameters:
//   - configName: The name of the configuration subdirectory where FFmpeg should be located.
//     This directory is created under the user's config directory (e.g., ~/.config/configName on Linux).
//
// # Returns:
//   - string: The full path to the FFmpeg binary if found and executable, empty string otherwise.
//   - bool: true if FFmpeg is found and executable in the specified location, false otherwise.
//
// # Example:
//
//	path, installed := IsStaticallyInstalled("my-app")
//	if installed {
//	    // Use the statically installed FFmpeg
//	    cmd := exec.Command(path, "-i", "input.mp4", "output.mp4")
//	} else {
//	    // FFmpeg not found in static location
//	    fmt.Println("FFmpeg not statically installed")
//	}
func IsStaticallyInstalled(configName string) (string, bool) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", false
	}

	binaryName := "ffmpeg"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	fullConfigDir := filepath.Join(configDir, configName)
	fullPath := filepath.Join(fullConfigDir, binaryName)
	_, err = exec.Command(fullPath, "-version").CombinedOutput()

	if err != nil {
		return "", false
	}

	return fullPath, true
}

// Download downloads and installs FFmpeg to the specified configuration directory.
//
// # Parameters:
//   - configName: The name of the configuration subdirectory where FFmpeg will be installed.
//     This directory is created under the user's config directory (e.g., ~/.config/configName on Linux).
//
// # Returns:
//   - string: The full path to the downloaded FFmpeg binary if successful.
//   - error: An error if any step of the download or installation process fails, including network errors, file system
//     errors, or extraction failures.
//
// # Example:
//
//	path, err := Download("my-app")
//	if err != nil {
//	    log.Fatalf("Failed to download FFmpeg: %v", err)
//	}
//
//	// Use the downloaded FFmpeg binary
//	cmd := exec.Command(path, "-i", "input.mp4", "output.mp4")
func Download(configName string) (string, error) {
	version := getLatestVersion()
	fileName := fmt.Sprintf("ffmpeg_%s_%s.zip", runtime.GOOS, runtime.GOARCH)
	url := fmt.Sprintf("https://github.com/vegidio/ffmpeg-downloader/releases/download/%s/%s", version, fileName)

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error getting user config directory: %v", err)
	}

	fullConfigDir := filepath.Join(configDir, configName)
	zipPath := filepath.Join(fullConfigDir, "ffmpeg.zip")

	err = download(url, zipPath)
	if err != nil {
		return "", fmt.Errorf("error downloading file: %v", err)
	}

	// Delete the download in the end
	defer os.Remove(zipPath)

	err = unzip(zipPath, fullConfigDir)
	if err != nil {
		return "", fmt.Errorf("error unziping file: %v", err)
	}

	binaryName := "ffmpeg"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	binaryPath := filepath.Join(fullConfigDir, binaryName)
	return binaryPath, nil
}

// GetFFmpegPath retrieves the path to the FFmpeg binary, ensuring it's available for use.
//
// This function implements a fallback strategy to locate or obtain FFmpeg:
//
//  1. First checks if FFmpeg is installed system-wide
//  2. If not found system-wide, checks for a statically installed version in the config directory
//  3. If neither is found, downloads and installs FFmpeg to the config directory
//
// When FFmpeg is found system-wide, it returns an empty string since the binary is accessible
// via the system PATH. For statically installed or downloaded versions, it returns the full
// path to the binary.
//
// # Parameters:
//   - configName: The name of the configuration directory where FFmpeg will be stored
//     if downloaded or where it should be looked for if statically installed.
//
// # Returns:
//   - string: The full path to the FFmpeg binary if statically installed or downloaded,
//     or an empty string if FFmpeg is available system-wide. Returns empty string
//     if download fails.
//
// # Example:
//
//	path := GetFFmpegPath("my-app")
//	if path == "" {
//	    // FFmpeg is available system-wide, use "ffmpeg" command directly
//	    cmd := exec.Command("ffmpeg", "-version")
//	} else {
//	    // Use the specific path returned
//	    cmd := exec.Command(path, "-version")
//	}
func GetFFmpegPath(configName string) string {
	installed := IsSystemInstalled()
	if installed {
		return ""
	}

	path, installed := IsStaticallyInstalled(configName)
	if installed {
		return path
	}

	path, err := Download(configName)
	if err != nil {
		return ""
	}

	return path
}
