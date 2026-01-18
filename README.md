# Godot Downloader

## Installation

```
go install github.com/xlanstar/godot-downloader
```

## Usage

- Download latest stable current platform version

```sh
godot-downloader
```

or equivalent

```sh
godot-downloader -v latest -p system
```

- Download latest experimental current platform version, and unzip it

```sh
godot-downloader -v latest-experimental -u
```

- Download v4.6-rc1 MacOS .Net version, and unzip it

```sh
godot-downloader -v 4.6-rc1 -p macos -u --mono
```
