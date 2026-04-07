package internal

import (
	"fmt"
	"net/http"
)

type Preset interface{
	Transform(input string)([]byte, error)
}

func route(name string)(Preset, error){
	switch name{
	case "S4":
	case "Seasonr4":
		return &preset.Race{}, nil
	case "mentor":
		return &preset.Simple{File: raceJSON}, nil
	default:
		return nil, fmt.Errorf("unkown preset: %s", name)
	}
}

func Season4()(http.Request){
	http.Request r();

	return r
}