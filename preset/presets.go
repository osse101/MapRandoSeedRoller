package preset

import (
	"embed"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
)

const DevVersion = 121

var TemplateMap = map[string]string{
	"s2":           "Community_Race_Season_2.json",
	"s3":           "Community_Race_Season_3_(No_animals).json",
	"s3a":          "Community_Race_Season_3_(Save_the_animals).json",
	"s4":           "Community_Race_Season_4.json",
	"default":      "Default.json",
	"mentor":       "Mentor_Tournament.json",
	"objectives":   "Winter_Tournament_-_4_Random_Objectives.json",
	"suits":        "Winter_Tournament_-_Double_Suit.json",
	"g91":          "Winter_Tournament_-_Gravity_9_+_1.json",
	"draft":        "Winter_Tournament_-_Item_Draft.json",
	"metroids":     "Winter_Tournament_-_Metroid_Objectives.json",
	"noobjectives": "Winter_Tournament_-_No_Objectives.json",
	"vmove":        "Winter_Tournament_-_Varia_+_Movement.json",
}

//go:embed data/*.json
var templateFS embed.FS

func LoadTemplate(name string) (map[string]interface{}, error) {
	fileName, ok := TemplateMap[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("preset '%s' not found", name)
	}
	template, err := templateFS.ReadFile("data/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading preset; %w", err)
	}
	var jsonObject map[string]interface{}
	err = json.Unmarshal(template, &jsonObject)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling preset; %w", err)
	}
	return jsonObject, nil
}

func GetPresetNames() []string {
	keys := slices.Sorted(maps.Keys(TemplateMap))
	return keys
}
