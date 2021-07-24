#include <EtherCard.h>
#include <SoftwareSerial.h>
#include "config.h"
#include "serial.h"
#include "ucma.h"

void setup() {
  setupSerial();
  softSerial.println(F("\n\n  --== setup started ==--"));
  setupEthernet();
  ucma::setup();
  softSerial.println(F("  --== setup finished ==--"));
}

void loop () {
  ucma::read(1, 0, data_t::performance);
  ether.packetLoop(ether.packetReceive());
}
