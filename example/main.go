package main

import (
	"fmt"
	ffmpegDownloader "github.com/vegidio/ffmpeg-downloader"
)

func main() {
	// Check if you already have FFmpeg in your computer
	if ffmpegDownloader.IsSystemInstalled() {
		fmt.Println("You already have FFmpeg installed in your system; no need to download anything")
		return
	}

	// Check if FFmpeg was previously installed using ffmpeg-downloader
	// In the code below, it looks for FFmpeg in the directory `<user-config>/ffmpeg-downloader`
	path, installed := ffmpegDownloader.IsStaticallyInstalled("ffmpeg-downloader")
	if installed {
		fmt.Println("You already downloaded FFmpeg using the ffmpeg-downlaoder; nothing to do here")
		fmt.Println("FFmpeg executable is installed in the path:", path)
		return
	}

	// Download the latest version of FFmpeg according to your computer's OS and architecture
	// In the code below, it downloads FFmpeg in the directory `<user-config>/ffmpeg-downloader`
	path, err := ffmpegDownloader.Download("ffmpeg-downloader")
	if err != nil {
		fmt.Println("Failed to download FFmpeg; see details:", err.Error())
		return
	}

	// Download successful. FFmpeg is now in the path `<user-config>/ffmpeg-downloader/ffmpeg`
	fmt.Println("Download successful; FFmpeg can now be found in the path:", path)
}
