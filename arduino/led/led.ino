#include <FastLED.h>

#define DATA_PIN    7
#define BRIGHTNESS  96
#define LED_TYPE    WS2812B
#define COLOR_ORDER GRB
#define NUM_LEDS    60 
CRGB leds[NUM_LEDS];

void setup() {
 delay(3000); // 3 second delay for recovery
  
  // tell FastLED about the LED strip configuration
  FastLED.addLeds<LED_TYPE,DATA_PIN,COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);

  // set master brightness control
  FastLED.setBrightness(BRIGHTNESS);
  
  //pinMode(LED_BUILTIN, OUTPUT);
  Serial.begin(9600);
}

void loop() { 
    EVERY_N_MILLISECONDS( 1 ) { 
      if(Serial.available() >= 4){
        byte start = Serial.read();
        if (start == 0) { // Black colour indicates next colour will be defined
          byte r = Serial.read();
          byte g = Serial.read();
          byte b = Serial.read();

          for( int i = 0; i < NUM_LEDS; i++) {
            leds[i] = CRGB(r, g, b);
          }
          // Send the 'leds' array out to the actual LED strip
          FastLED.show();  
        }
      }
    }
}
