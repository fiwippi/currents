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
      if(Serial.available()){
        int num = Serial.parseInt();
        if (num == -1) { // Indicates next colour will be defined
          unsigned long hue = Serial.parseInt();
          int r = hue & 0xFF;
          int g = (hue >> 8) & 0xFF;
          int b = (hue >> 16) & 0xFF;

          for( int i = 0; i < NUM_LEDS; i++) {
            leds[i] = CRGB(r, g, b);
          }
          // Send the 'leds' array out to the actual LED strip
          FastLED.show();  
        }
      }
    }
}
