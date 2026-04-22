package lib

import (
	"os"
)

func BuildSpoilerToken() string {
	return os.Getenv("SPOILER_TOKEN")
}

func BuildSite(isDev bool) string {
	if isDev {
		return "https://dev.maprando.com"
	}
	return "https://maprando.com"
}
