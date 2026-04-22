package lib

import (
	"os"
	"testing"
)

func TestBuildSpoilerToken(t *testing.T) {
	os.Setenv("SPOILER_TOKEN", "test_token")
	defer os.Unsetenv("SPOILER_TOKEN")

	token := BuildSpoilerToken()
	if token != "test_token" {
		t.Errorf("Expected token 'test_token', got '%s'", token)
	}
}

func TestBuildSite(t *testing.T) {
	devSite := BuildSite(true)
	if devSite != "https://dev.maprando.com" {
		t.Errorf("Expected 'https://dev.maprando.com', got '%s'", devSite)
	}

	prodSite := BuildSite(false)
	if prodSite != "https://maprando.com" {
		t.Errorf("Expected 'https://maprando.com', got '%s'", prodSite)
	}
}
