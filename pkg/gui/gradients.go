package gui

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/lucasb-eyer/go-colorful"

	"currents/internal/log"
	"currents/pkg/audio"
)

func loadGradients() *audio.Gradients {
	// Create hardcoded gradients
	gradients := audio.NewGradients()
	gradients.Add("Piercing", audio.Gradient{
		{audio.MustParseHex("#1152cb"), 0.0},
		{audio.MustParseHex("#e4032f"), 0.63869},
	},
	)
	gradients.Add("Warm", audio.Gradient{
		{audio.MustParseHex("#ff0000"), 0},
		{colorful.Color{R: 1, G: 0.8627450980392157, B: 0}, 0.6993007063865662},
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
