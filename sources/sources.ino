#include <EtherCard.h>

void setup() {
  setupEthernet();
}

void loop () {
  ether.packetLoop(ether.packetReceive());
}
