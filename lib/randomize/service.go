package randomize

import (
	"maprandoseedroller/lib"
	"maprandoseedroller/lib/models"
)

func Randomize(data []byte, isDev bool) (string, error) {
	//Make request to Map Rando
	r := models.RequestMapRando{
		Settings:     data,
		SpoilerToken: lib.BuildSpoilerToken(),
	}

	resp, err := MakeRequest(lib.BuildSite((isDev)), r)
	if err != nil {
		return "", err
	}

	//Return result
	return resp, nil
}
