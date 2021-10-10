package audio

type Config struct {
	Channels   uint32 // Number of channels, default is 2
	SampleRate uint32 // Sample rate, default is 44100
}

func DefaultConfig() *Config {
	return &Config{
		Channels:   2,
		SampleRate: 44100,
	}
}
