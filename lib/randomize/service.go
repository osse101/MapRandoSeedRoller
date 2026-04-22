package randomize

import (
	"maprandoseedroller/lib/models"
)

func Randomize(data []byte)(string, error){
	//Make request to Map Rando
	r := models.RequestMapRando{
		Settings: data,
		SpoilerToken: "s",
	}

	resp, err := RandomizeRequest("maprando.com/", r)
	if err != nil{
		return "", err
	}

	//Return result
	return resp, nil
}