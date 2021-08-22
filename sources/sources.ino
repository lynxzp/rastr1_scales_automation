#include <EtherCard.h>
#include <SoftwareSerial.h>
#include "config.h"
#include "serial.h"
#include "ucma.h"
#include "ethernet.h"
#include "modbus.h"

void setup() {
  setupSerial();
  softSerial.println(F("\n\n  --== setup started ==--"));
  setupEthernet();
  ucma::setup();
  softSerial.println(F("  --== setup finished ==--"));
  softSerial.print(F("  --== start read ucma ==--\n"));
}

static BufferFiller bfill;
Modbus modbus;

void loop () {
    /*int32_t data = ucma::read(2, data_t::accumulation);
    softSerial.print("accumulation: ");
    char buf[10];
    sprintf(buf, "%7ld", data);
    softSerial.println(buf);
    delay(500);
    data = ucma::read(2, data_t::performance2avg);
    softSerial.print("performance:  ");
    char buf[10];
    sprintf(buf, "%5d.", data/10);
    softSerial.print(buf);
    softSerial.println(data%10);
    delay(5000);
    softSerial.println();*/

    uint16_t payloadPos = ether.packetLoop(ether.packetReceive());
    if (payloadPos)
    {
        char* incomingData = (char *) Ethernet::buffer + payloadPos;
        if(modbus.encodeTCP(incomingData))
        {
            auto uintId = modbus.getUnitIdentifier();
            auto dataAddr = modbus.getRequestedDataAddress();
            dataAddr = 0x60;    // !!!! remove
            uintId = 2;     // !!!! remove
            auto data = ucma::read(uintId, (data_t)(dataAddr));
            softSerial.print(F("read: "));
            softSerial.println(data);
            modbus.setData(data);

            bfill = ether.tcpOffset();
            bfill.emit_raw(modbus.getResponseBuf(), modbus.getResponseSize());
        }
    }

}
