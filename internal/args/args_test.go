package args_test

import (
	"godot-downloader/internal/args"
	"testing"
)

func TestIsVersionValid(t *testing.T) {
	testCases := []struct {
		version string
		want    bool
	}{
		{"1.2", true},
		{"1.2.3", true},
		{"12.347.5689", true},
		{"4.1", true},
		{"4.5", true},
		{"3.6", true},
		{"1.2.3-alpha", false},
		{"1.2.3-beta", false},
		{"1.2.3-rc", false},
		{"invalid", false},
	}

	for _, tc := range testCases {
		got := args.IsVersionValid(tc.version)
		if got != tc.want {
			t.Errorf("testVersion(%q) = %v, want %v", tc.version, got, tc.want)
		}
	}
}

func TestIsSlugValid(t *testing.T) {
	testCases := []struct {
		slug string
		want bool
	}{
		{"stable", true},
		{"dev1", true},
		{"alpha2", true},
		{"beta3", true},
		{"rc4", true},
		{"invalid", false},
		{"dev", false},
		{"alpha", false},
		{"beta", false},
		{"rc", false},
	}

	for _, tc := range testCases {
		got := args.IsSlugValid(tc.slug)
		if got != tc.want {
			t.Errorf("testSlug(%q) = %v, want %v", tc.slug, got, tc.want)
		}
	}
}
