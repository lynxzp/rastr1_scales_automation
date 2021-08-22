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
const uint8_t nums_len = 4;
uint8_t nums[nums_len] {55,63,68,93};
void loop () {
//    while(1) {
//        for(int j=0; j<nums_len; j++) {
//            auto data = ucma::read(2, data_t(nums[j]));
//            softSerial.print(nums[j]);
//            softSerial.print(" | ");
////            softSerial.print(data);
////            if((data>20000)&&(data<60000)){
////                softSerial.print(" !");
////            }
////            softSerial.println();
////            delay(100);
//        }
//        softSerial.println();
//        delay(1000);
//    }
    int32_t data = ucma::read(2, data_t::accumulation);
    /*if(data!=-1)*/ {
      softSerial.print("accumulation: ");
      char buf[10];
      sprintf(buf, "%7ld", data);
      softSerial.println(buf);
    }
    delay(500);
    data = ucma::read(2, data_t::performance2avg);
    /*if(data!=-1)*/ {
        softSerial.print("performance:  ");
        char buf[10];
        sprintf(buf, "%5d.", data/10);
        softSerial.print(buf);
        softSerial.println(data%10);
    }
    delay(5000);
    softSerial.println();
//    ether.packetLoop(ether.packetReceive());
}
