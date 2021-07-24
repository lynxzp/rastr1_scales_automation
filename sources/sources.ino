#include <EtherCard.h>
#include <SoftwareSerial.h>
#include "config.h"
#include "serial.h"

void setup() {
  setupSerial();
  softSerial.println(F("\n\n  --== setupEthernet ==--"));
  setupEthernet();
}

void loop () {
  ether.packetLoop(ether.packetReceive());
}
