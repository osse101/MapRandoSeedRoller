package api

import (
	"fmt"
	"net/http"
)

func HelpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a preset name.")
	return
}