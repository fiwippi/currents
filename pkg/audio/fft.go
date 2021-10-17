package audio

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"math/cmplx"
	"time"

	"github.com/lucasb-eyer/go-colorful"

	"currents/internal/fft"
)

var ErrChannelNum = errors.New("channel number is unsupported")
var ErrSampleRate = errors.New("sample rate is unsupported")

type FFT struct {
	buffer    *bytes.Buffer // Writer the captured audio should be written to
	abortChan chan error    // Tells FFT to stop processing input on its writer
	conf      *Config       // Config tells FFT how to process the audio data

	// What type of interpolation to use for drawing the colours
	DrawMode InterpolateMode
	// Chan which returns the most recent calculated colour
	Hues chan colorful.Color
	// When FFT has stopped processing, due to error or not,
	// the error message is sent here
	Done chan error
	// Gradient to interpolate colours with, if this is not specified
	// then colours are interpolated over the HSV spectrum
	Gradient *Gradient
	// Whether the hue colour change should be damped
	Damp bool
	// How many past frequencies are taken into account when damping
	DampSliceLen int
	// Whether the hue colour change should be smoothed
	Smooth bool
	// How strong should the smoothing be from 0 <= alpha < 1,
	// at other values the smoothing behaviour is undefined
	SmoothAlpha float64
	// FFT will clamp the maximum frequency to this value
	MaxFreq float64
	// The upper range of frequencies the program considers useful.
	// After this barrier, the colour changes very slowly in relation
	// to change in f
	MaxUsefulFrequency float64
	// How many unique shoes should the colour spectrum have
	TotalHues float64
	// The hue number at which the MaxUsefulFrequency is reached,
	// e.g. from 200 total hues, after 180 hues MaxUsefulFrequency
	// will take effect
	UsefulFrequencyHue float64
	// How often we want to use the values from the audio buffer
	SampleRate time.Duration
}

func NewFFT(conf *Config) (*FFT, error) {
	if conf.Channels != 2 {
		return nil, ErrChannelNum
	} else if conf.SampleRate != 44100 {
		return nil, ErrSampleRate
	}

	f := &FFT{
		buffer:             bytes.NewBuffer(nil),
		abortChan:          make(chan error),
		conf:               conf,
		DrawMode:           Blended,
		Hues:               make(chan colorful.Color, 1),
		Done:               make(chan error),
		MaxFreq:            2500,
		MaxUsefulFrequency: 1200,
		TotalHues:          320,
		UsefulFrequencyHue: 310,
		Damp:               true,
		DampSliceLen:       4,
		Smooth:             true,
		SmoothAlpha:        0.65,
		SampleRate:         75 * time.Millisecond,
	}
	go f.start()
	return f, nil
}

func MustCreateNewFFT(conf *Config) *FFT {
	f, err := NewFFT(conf)
	if err != nil {
		panic(err)
	}
	return f
}

// Write implements io.Writer
func (f *FFT) Write(p []byte) (n int, err error) {
	return f.buffer.Write(p)
}

func (f *FFT) start() {
	channelNum := int(f.conf.Channels)
	sampleRate := int(f.conf.SampleRate)

	buf := make([]byte, 0, 1024*8)
	reader := bufio.NewReader(f.buffer)

	var displayFreq float64 // Interpolated frequency displayed on the LED lights
	var frequency float64   // The max frequency of the current buffer
	var oldFreq float64     // The max frequency of the previous buffer

	rate := time.NewTicker(f.SampleRate)
	var update time.Time

	for {
		select {
		case err := <-f.abortChan:
			close(f.Hues)
			rate.Stop()
			f.Done <- err
			return
		default:
			// Populate the buffer
			n, err := fill(reader, buf)
			buf = buf[:n]
			if err != nil && err != io.ErrUnexpectedEOF {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					continue
				}
				close(f.Hues)
				f.Done <- err
				return
			}

			// Frequency only updated every delta t, colour
			// updated instantaneously
			select {
			case <-rate.C:
				update = time.Now()
				// Get all the float values for each sample in the input samples
				monoFrameCount := len(buf) / (channelNum * sampleSizeInBytes) // We mix down the samples into mono so we lose half the frames
				samples := make([]float32, 0, monoFrameCount)
				for i := 0; i < len(buf); i += channelNum * sampleSizeInBytes {
					// Value of th left sample
					leftBytes := buf[i : i+sampleSizeInBytes]
					leftBits := binary.LittleEndian.Uint32(leftBytes)
					leftFloat := math.Float32frombits(leftBits)

					// Value of the right sample
					rightBytes := buf[i+sampleSizeInBytes : i+channelNum*sampleSizeInBytes]
					rightBits := binary.LittleEndian.Uint32(rightBytes)
					rightFloat := math.Float32frombits(rightBits)

					// Mix them together
					mixedFloat := (leftFloat + rightFloat) / float32(channelNum)
					samples = append(samples, mixedFloat)
				}

				// Calculate variables used to get the FFT
				maxInfo := sampleRate / channelNum                  // Reversing the Nyquist–Shannon sampling theorem to see the maximum frequency we are trying to achieve
				usefulMonoFrameCount := monoFrameCount / channelNum // This is the length the program uses to find the freq with
				// the highest magnitude, this is half the buffer length because
				// the FFT is mirrored along the centre, thus only half the length
				// needs to be used
				freqBinSize := maxInfo / usefulMonoFrameCount // This represents the difference in frequency between each index of the FFT'd array

				// Perform the FFT on the samples and get the frequency with the largest magnitude
				fftData := fft.FFTReal(samples)

				var max float64
				var index int

				// FFT is mirrored so we only need the first half of the samples
				for i := 0; i < usefulMonoFrameCount; i++ {
					e := cmplx.Abs(complex128(fftData[i]))
					if e > max {
						max = e
						index = i
					}
				}

				oldFreq = frequency
				frequency = math.Min(float64(freqBinSize*index), f.MaxFreq)

				// Smooth if needed
				if f.Smooth {
					frequency = (f.SmoothAlpha * oldFreq) + ((1 - f.SmoothAlpha) * frequency)
				}
			default:
			}

			displayFreq = frequency

			// Damp if needed
			if f.Damp {
				since := time.Now().Sub(update)
				delta := float64(since.Nanoseconds()) / float64(f.SampleRate.Nanoseconds())
				displayFreq = oldFreq + delta*(frequency-oldFreq)
			}

			// Calculate the corresponding hue for the colour
			var hue float64
			if displayFreq > f.MaxUsefulFrequency {
				hue = f.UsefulFrequencyHue + (f.TotalHues-f.UsefulFrequencyHue)*(displayFreq/f.MaxFreq)
			} else {
				hue = displayFreq / f.MaxUsefulFrequency * f.UsefulFrequencyHue
			}

			// Create the colour
			var colour colorful.Color
			if f.Gradient != nil {
				colour = f.DrawMode.Interpolate(hue/f.TotalHues, *f.Gradient)
			} else {
				colour = colorful.Hsv(hue, 1, 1)
			}
			f.Hues <- colour
		}
	}
}

func (f *FFT) Stop() {
	f.abortChan <- nil
}

// fill ignores io.EOF and io.ErrUnexpectedEOF and waits until the buffer
// is full before returning. This is useful because a lot of times the fft
// buffer will be read from at around the same speed it's written to so it
// can't fill up completely to its full capacity, this fixes that
func fill(reader io.Reader, buf []byte) (int, error) {
	totalN := 0

	for {
		n, err := io.ReadFull(reader, buf[totalN:cap(buf)])
		totalN += n
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return totalN, err
		}

		if cap(buf)-totalN == 0 {
			return totalN, nil
		}
	}
}
