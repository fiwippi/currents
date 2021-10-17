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
		{audio.MustParseHex("#1152cb"), 0.05},
		{audio.MustParseHex("#e4032f"), 0.17482},
		{audio.MustParseHex("#f6c507"), 0.83449},
		{audio.MustParseHex("#faf6cb"), 1.0},
	},
	)
	gradients.Add("Warm", audio.Gradient{
		{audio.MustParseHex("#ff0000"), 0},
		{colorful.Color{R: 0.42745098039215684, G: 0, B: 1}, 0.4265},
	},
	)
	gradients.Add("Cold", audio.Gradient{
		{colorful.Color{R: 0, G: 0.8666666666666667, B: 1}, 0},
		{colorful.Color{R: 0, G: 0.8666666666666667, B: 1}, 1},
	},
	)
	gradients.Add("Jungle", audio.Gradient{
		{colorful.Color{R: 0, G: 1, B: 0}, 0},
		{colorful.Color{R: 0.7450980392156863, G: 1, B: 1}, 0.16317015886306763},
		{colorful.Color{R: 0.7333333333333333, G: 0.3686274509803922, B: 0.3686274509803922}, 0.7785547971725464},
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
