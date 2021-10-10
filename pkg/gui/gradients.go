package gui

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"currents/internal/log"
	"currents/pkg/audio"
)

func loadGradients() *audio.Gradients {
	// Create hardcoded gradients
	gradients := audio.NewGradients()
	gradients.Add("Starboy", audio.Gradient{
		{audio.MustParseHex("#1152cb"), 0.0},
		{audio.MustParseHex("#1152cb"), 0.05},
		{audio.MustParseHex("#e4032f"), 0.1},
		{audio.MustParseHex("#f6c507"), 0.55},
		{audio.MustParseHex("#faf6cb"), 1.0},
	},
	)
	gradients.Add("Smiths", audio.Gradient{
		{audio.MustParseHex("#ff0202"), 0},
		{audio.MustParseHex("#ff0202"), 0.1},
		{audio.MustParseHex("#ff8d00"), 0.3},
		{audio.MustParseHex("#fff400"), 0.5},
		{audio.MustParseHex("#f1ff00"), 0.8},
		{audio.MustParseHex("#A4ff00"), 1.0},
	},
	)
	gradients.Add("Weird", audio.Gradient{
		{audio.MustParseHex("#ff7303"), 0},
		{audio.MustParseHex("#ff7303"), 0.1},
		{audio.MustParseHex("#ffa7e1"), 0.5},
		{audio.MustParseHex("#faf4e6"), 1.0},
	},
	)

	// Add custom gradients if any exist, they are allowed to overwrite the custom ones
	f, err := os.OpenFile("gradients.json", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("could not open gradients.json")
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read gradients.json")
	}

	if len(data) > 0 {
		var customGradients audio.Gradients
		err = json.Unmarshal(data, &customGradients)
		if err != nil {
			log.Fatal().Err(err).Msg("could not unmarshal from gradients.json")
		}

		for _, name := range customGradients.List() {
			gradients.Add(name, customGradients.Get(name))
		}
	}

	return gradients
}
