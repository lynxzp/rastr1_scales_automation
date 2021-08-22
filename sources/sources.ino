#include <EtherCard.h>
#include <SoftwareSerial.h>
#include "config.h"
#include "serial.h"
#include "ucma.h"
#include "ethernet.h"

void setup() {
  setupSerial();
  softSerial.println(F("\n\n  --== setup started ==--"));
  //setupEthernet();
  ucma::setup();
  softSerial.println(F("  --== setup finished ==--"));
  softSerial.print(F("  --== start read ucma ==--\n"));
}

void push(uint8_t *buf,uint8_t len) {
  for(uint8_t i=0;i<len;i++) {
    Serial.print(char(buf[i]));
  }
}

bool serialread() {
  if(Serial.available()) {
    softSerial.print("received:");
    softSerial.println(int(Serial.read()));
    return true;
  }
  return false;
}

void printbuf(uint8_t* buf, uint8_t len){
  for(int i=0;i<len;i++) {
    softSerial.print(int(buf[i]));
    softSerial.print(",");
    
  }
  softSerial.println("");
}

void loop () {
    int32_t data = ucma::read(2, data_t::accumulation);
    /*if(data!=-1)*/ {
      softSerial.print("accumulation:");
      softSerial.println(data);
    }
    delay(500);
    data = ucma::read(2, data_t::performance);
    /*if(data!=-1)*/ {
        softSerial.print("performance:");
        softSerial.println(data);
    }
    delay(5000);
    softSerial.println();
//    ether.packetLoop(ether.packetReceive());
}
