# ffmpeg-downloader

A Go library to automatically download FFmpeg in your project.

## ⬇️ Installation

This library can be installed using Go modules. To do that run the following command in your project's root directory:

```bash
$ go get github.com/vegidio/ffmpeg-downloader
```

## 🤖 Usage

**ffpmeg-downloader** exposes three functions to check the status of FFmpeg and manage its download:

- `IsSystemInstalled()`: Checks whether FFmpeg is already installed on the system by the user. If it returns `true`, you likely don't need to download FFmpeg, as you can use the version available in the system's `PATH`.
- `IsStaticallyInstalled(<name>)`: Checks whether FFmpeg was previously installed using **ffmpeg-downloader**. It returns the installation path and `true` if found; otherwise, it returns an empty string and `false`.
- `Download(<name>)`: Downloads the latest version of FFmpeg to the system. If successful, it returns the installation path and `nil`; otherwise, it returns an empty string and an `error`.

You can find a working example in the [example](example) folder.

### Where is FFmpeg installed?

The `IsStaticallyInstalled(<name>)` and `Download(<name>)` functions use the user's configuration directory to locate or install FFmpeg. The path to this config directory depends on the operating system:

- **Linux**: `/home/<username>/.config/<name>`
- **macOS**: `/Users/<username>/Library/Application Support/<name>`
- **Windows**: `C:\Users\<username>\AppData\Roaming\<name>`

Here, `<username>` refers to the current user running the application, and `<name>` is the parameter you provide to the functions - typically your app's name. It’s essential to use the same `<name>` across both functions to ensure the library can correctly locate the FFmpeg installation.

## 📝 License

**ffmpeg-downloader** is released under the MIT License. See [LICENSE](LICENSE) for details.

## 👨🏾‍💻 Author

Vinicius Egidio ([vinicius.io](http://vinicius.io))