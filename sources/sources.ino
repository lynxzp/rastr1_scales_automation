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
  softSerial.println(F("  --== start read ucma ==--"));
}

static BufferFiller bfill;
Modbus modbus;


static word homePage() {
  long t = millis() / 1000;
  word h = t / 3600;
  byte m = (t / 60) % 60;
  byte s = t % 60;
  bfill = ether.tcpOffset();
  bfill.emit_p(PSTR("$D$D:$D$D:$D$D"),
      h/10, h%10, m/10, m%10, s/10, s%10);
  return bfill.position();
}


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

    /*uint16_t payloadPos = ether.packetLoop(ether.packetReceive());
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
            char str[10]="Hello 123";
            bfill.emit_p(PSTR("==$F"),str);
            //bfill.emit_raw("123abc\n",7);
            //bfill.emit_raw(modbus.getResponseBuf(), modbus.getResponseSize());
            softSerial.print(F("was sent\n"));
        }
    }*/
  word pos = ether.packetLoop(ether.packetReceive());
  if (pos){
    char* incomingData = (char *) Ethernet::buffer + pos;
    if(modbus.encodeTCP(incomingData)){
      ether.httpServerReply(homePage()); // send web page data
    }
  }

    auto now_time = millis();
    static uint32_t last_processed_time = now_time;
    if(now_time>last_processed_time+10000) {
        last_processed_time+=10000;
        softSerial.println(last_processed_time/10000);
    }
}
