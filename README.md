![currents](icon.png)
# Currents
Currents changes the colour of your LED lights based on music playing on your device

## Requirements
Your audio device must:
- 2 Channels
- 441000 Hz Sample Rate
- F32 Sound (as specified by Mini Audio)

## Usage
### Arduino
1. Build an Arduino according to the specification [here](https://www.makeuseof.com/tag/connect-led-light-strips-arduino/) and use a WS2812B LED strip
2. Load the program from ./arduino/led/led.ino onto your arduino (modify it if you need to change data pin/led type/number of leds)

### GUI
1. Clone the repo using `git clone https://github.com/fiwippi/currents.git`
2. Build the binary `cd currents && make build`
3. Send the audio you want to visualise to a virtual audio cable (so it can act as a mic)
4. Select this virtual audio cable when running currents

## License
`BSD-3-Clause`