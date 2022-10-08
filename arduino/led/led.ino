#include <FastLED.h>

#define DATA_PIN    7
#define BRIGHTNESS  120
#define LED_TYPE    WS2812B
#define COLOR_ORDER GRB
#define NUM_LEDS    300

uint8_t buf[8]; 
CRGB leds[NUM_LEDS];

char endSeq[] = {0xCC, 0xDD};

void setup() {
  delay(3000); // 3 second delay for recovery
  
  // tell FastLED about the LED strip configuration
  FastLED.addLeds<LED_TYPE,DATA_PIN,COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);

  // set master brightness control
  FastLED.setBrightness(BRIGHTNESS);

  Serial.setTimeout(5);
  Serial.begin(9600);
}

void loop() { 
      if(Serial.available()){
        int rlen = Serial.readBytes(buf, 8);
        if (rlen == 8) {
          Serial.write(buf[0]);

          int valid = 0;
          CRGB clr;

          if (buf[0] == 0xAA && buf[1] == 0xBB && buf[2] == 0x00 && buf[6] == 0xCC && buf[7] == 0xDD ) {
            clr = CRGB(buf[3], buf[4], buf[5]);
            valid = 1;
          } else if (buf[0] == 0xBB && buf[1] == 0x00 && buf[5] == 0xCC && buf[6] == 0xDD ) {
            clr = CRGB(buf[2], buf[3], buf[4]);
            valid = 1;
          } else if (buf[0] == 0x00 && buf[4] == 0xCC && buf[5] == 0xDD ) {
            clr = CRGB(buf[1], buf[2], buf[3]);
            valid = 1;
          }
          
          if (valid) {
            for( int i = 0; i < NUM_LEDS; i++) {
              leds[i] = clr;
            }
            FastLED.show();  
          } else if (buf[0] == 0xCC || buf[0] == 0xDD) {
            Serial.write(255);
            while (!Serial.find(endSeq)) {Serial.write(254);};
          }
        }
      }
}
