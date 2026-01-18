package args

import (
	"fmt"
	"godot-downloader/internal/parser"
	"regexp"
	"runtime"
	"strings"
)

type GodotPlatformArg string

const (
	GodotLatestVersion             = "latest"
	GodotLatestExperimentalVersion = "latest-experimental"
)

const (
	MacOS   GodotPlatformArg = "macos"
	Linux   GodotPlatformArg = "linux"
	Windows GodotPlatformArg = "windows"
	Android GodotPlatformArg = "android"
	Web     GodotPlatformArg = "web"
	System  GodotPlatformArg = "system"
)

var Args struct {
	// Specific Godot Version
	// 4.4-stable, 4.5-dev1 4.5-beta2 4.5-rc3
	// "latest" means the latest minor stable version
	// "latest-experimental" means the latest minor experimental version
	GodotVersion string `arg:"-v" default:"latest"`

	// Specific Platform
	// macOS, Linux, Windows, Android, Web.
	// "system" means the current platform
	GodotPlatform GodotPlatformArg `arg:"-p" default:"system"`

	// Specific Mono version
	GodotMono bool `arg:"--mono" default:"false"`

	// Specific Output File Path
	OutputFilePath string `arg:"-o"`

	// Unarchive the downloaded file
	Unarchive bool `arg:"-u" default:"false"`
}

func ResolveArgs() {
	switch Args.GodotVersion {
	case GodotLatestVersion:
		Args.GodotVersion = parser.GetLatestVersion()
	case GodotLatestExperimentalVersion:
		Args.GodotVersion = parser.GetLatestExperimentalVersion()
	}

	switch Args.GodotPlatform {
	case System:
		switch runtime.GOOS {
		case "windows":
			Args.GodotPlatform = Windows
		case "darwin":
			Args.GodotPlatform = MacOS
		case "linux":
			Args.GodotPlatform = Linux
		case "android":
			Args.GodotPlatform = Android
		case "web":
			Args.GodotPlatform = Web
		default:
			fmt.Printf("Unsupported or unknown OS: %s\n", runtime.GOOS)
		}
	}
}

func ParseGodotVersionAndSlug(godotVersionArg string) (version, slug string) {
	if !strings.Contains(godotVersionArg, "-") {
		return
	}

	segs := strings.Split(godotVersionArg, "-")
	version, slug = segs[0], segs[1]

	if !IsVersionValid(version) {
		return
	}

	if !IsSlugValid(slug) {
		return
	}

	return version, slug
}

func IsVersionValid(version string) bool {
	matched, _ := regexp.MatchString(`^[0-9]+\.[0-9]+(\.[0-9]+)?$`, version)
	return matched
}

func IsSlugValid(slug string) bool {
	matched, _ := regexp.MatchString(`^(stable|(alpha|dev|beta|rc)\d)$`, slug)
	return matched
}
