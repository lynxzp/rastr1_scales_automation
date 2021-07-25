#include <EtherCard.h>
#include <SoftwareSerial.h>
#include "config.h"
#include "serial.h"
#include "ucma.h"
#include "ethernet.h"

void setup() {
  setupSerial();
  softSerial.println(F("\n\n  --== setup started ==--"));
  setupEthernet();
  ucma::setup();
  softSerial.println(F("  --== setup finished ==--"));
  softSerial.print(F("  --== start read ucma ==--\n asking address:"));
}

int addr = 0;

void loop () {
  softSerial.print(addr);
  softSerial.print(" ");
  ucma::read(addr++, 0, data_t::accumulation);
  ether.packetLoop(ether.packetReceive());
  
}
