package ffmpeg_downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// IsSystemInstalled checks if FFmpeg is installed on the system.
//
// It returns true if it's installed, otherwise it returns false.
func IsSystemInstalled() bool {
	_, err := exec.Command("ffmpeg", "-version").CombinedOutput()
	if err != nil {
		return false
	}

	return true
}

// IsStaticallyInstalled checks if the FFmpeg binary is installed in the specified configuration directory.
//
// It takes a configuration name as an argument, which is used to determine the directory where the FFmpeg binary
// should be located.
//
// Parameters:
//   - configName: The name of the configuration directory.
//
// Returns:
//   - string: The path to the FFmpeg binary if found, otherwise an empty string.
//   - bool: true if the FFmpeg binary is found and executable, otherwise false.
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

// Download downloads and installs the FFmpeg binary.
//
// It takes a configuration name as an argument, which is used to determine the directory where the FFmpeg binary will
// be installed.
//
// Parameters:
//   - configName: The name of the configuration directory.
//
// Returns:
//   - string: The path to the installed FFmpeg binary.
//   - error: An error message if the installation fails.
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
