#include <EtherCard.h>
#include <SoftwareSerial.h>
#include "config.h"
#include "serial.h"
#include "ucma.h"
#include "ethernet.h"
#include "modbus.h"

void setup() {
  setupSerial();
  //softSerial.println(F("\n\n  --== setup started ==--"));
  setupEthernet();
  ucma::setup();
  //softSerial.println(F("  --== setup finished ==--"));
  //softSerial.println(F("  --== start read ucma ==--"));
}

static BufferFiller bfill;
Modbus modbus;


static word responsePage() {
  auto ptr = modbus.getResponseBuf();
  bfill = ether.tcpOffset();
  bfill.emit_p(PSTR("{\"transaction\":$D, \"unit\":$D, \"data\":$L}"),
    uint16_t(ptr->transactionIdentifier),
    uint16_t(ptr->unitIdentifier),
    int32_t(ptr->data));
    //softSerial.println(ptr->data);
  return bfill.position();
}


void loop () {
  word pos = ether.packetLoop(ether.packetReceive());
  if (pos){
    char* incomingData = (char *) Ethernet::buffer + pos;
    if(modbus.encodeTCP(incomingData)){
      auto uintId = modbus.getUnitIdentifier();
      auto dataAddr = modbus.getRequestedDataAddress();
      auto data = ucma::read(uintId, (data_t)(dataAddr));
      modbus.setData(data);
      ether.httpServerReply(responsePage());
    }
  }
}
