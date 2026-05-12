package workflow

import (
	"encoding/json"
	"fmt"

	"maprandoseedroller/lib/models"
)

func Process(req models.RequestRaw) (models.ResponseOut3, error) {
	if req.Action != "" {
		return handleAction(req)
	}
	if req.Event != "" {
		return handleEvent(req)
	}
	return models.ResponseOut3{Status: "error"}, fmt.Errorf("no action or event specified")
}

func handleAction(req models.RequestRaw) (models.ResponseOut3, error) {
	switch req.Action {
	case "roll":
		var flags string
		if err := json.Unmarshal(req.Data, &flags); err != nil {
			return models.ResponseOut3{Status: "error"}, fmt.Errorf("roll action requires string flags")
		}
		res, err := ExecuteRoll(flags)
		return models.ResponseOut3{Status: "success", Data: res}, err

	case "unlock":
		var url string
		if err := json.Unmarshal(req.Data, &url); err != nil {
			return models.ResponseOut3{Status: "error"}, fmt.Errorf("unlock action requires a URL string")
		}
		res, err := ExecuteUnlock(url)
		return models.ResponseOut3{Status: "success", Data: res}, err

	case "help":
		msg := GetHelp(req.Source)
		return models.ResponseOut3{Status: "success", Message: msg}, nil

	default:
		return models.ResponseOut3{Status: "error"}, fmt.Errorf("unrecognized action: %s", req.Action)
	}
}

func handleEvent(req models.RequestRaw) (models.ResponseOut3, error) {
	switch req.Event {
	case "seed.finished":
		return models.ResponseOut3{Status: "success", Message: "Log acknowledged"}, nil
	case "system.alert":
		return models.ResponseOut3{Status: "success"}, nil
	default:
		// Silently acknowledge unknown events to stop webhook retries
		return models.ResponseOut3{Status: "success", Message: "event ignored"}, nil
	}
}
