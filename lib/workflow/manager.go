package workflow

import (
	"encoding/json"
	"fmt"

	"maprandoseedroller/lib/models"
)

func Process(req models.RequestRaw) (models.ResponseOut, error) {
	if req.Action != "" {
		return handleAction(req)
	}
	if req.Event != "" {
		return handleEvent(req)
	}
	return models.ResponseOut{Status: models.StatusError}, fmt.Errorf("no action or event specified")
}

func handleAction(req models.RequestRaw) (models.ResponseOut, error) {
	switch req.Action {
	case "roll":
		var flags string
		if err := json.Unmarshal(req.Data, &flags); err != nil {
			return models.ResponseOut{Status: models.StatusError}, fmt.Errorf("roll action requires string flags")
		}
		res, err := ExecuteRoll(flags)
		if err != nil {
			return models.ResponseOut{Status: models.StatusError}, err
		}
		return models.ResponseOut{Status: models.StatusSuccess, Data: res}, nil

	case "unlock":
		var url string
		if err := json.Unmarshal(req.Data, &url); err != nil {
			return models.ResponseOut{Status: models.StatusError}, fmt.Errorf("unlock action requires a URL string")
		}
		if err := ExecuteUnlock(url); err != nil {
			return models.ResponseOut{Status: models.StatusError}, err
		}
		return models.ResponseOut{Status: models.StatusSuccess}, nil

	case "help":
		msg := GetHelp(req.Source)
		return models.ResponseOut{Status: models.StatusSuccess, Message: msg}, nil

	default:
		return models.ResponseOut{Status: models.StatusError}, fmt.Errorf("unrecognized action: %s", req.Action)
	}
}

func handleEvent(req models.RequestRaw) (models.ResponseOut, error) {
	switch req.Event {
	case "seed.finished":
		return models.ResponseOut{Status: models.StatusSuccess, Message: "Log acknowledged"}, nil
	case "system.alert":
		return models.ResponseOut{Status: models.StatusSuccess}, nil
	default:
		// Silently acknowledge unknown events to stop webhook retries
		return models.ResponseOut{Status: models.StatusSuccess, Message: "event ignored"}, nil
	}
}
