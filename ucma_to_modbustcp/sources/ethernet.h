uint8_t Ethernet::buffer[ethernetBufferSize];

void printIp (const uint8_t *buf) {
    for (uint8_t i = 0; i < IP_LEN; ++i) {
        softSerial.print( buf[i], DEC );
        if (i < 3)
            softSerial.print('.');
    }
    softSerial.println();
}

void setupEthernet() {
    ether.hisport = modbusTcpPort;
  softSerial.print("MAC: ");
  for (byte i = 0; i < 6; ++i) {
    softSerial.print(ethernetMAC[i], HEX);
    if (i < 5)
      softSerial.print(':');
  }
  softSerial.println();

  if (ether.begin(sizeof Ethernet::buffer, ethernetMAC, ethernetCS) == 0){
    softSerial.println(F("Failed to access Ethernet controller"));
  } else {
    softSerial.println(F("Sucessfull access Ethernet controller"));
  }

  softSerial.println(F("Setting up DHCP"));
  if (!ether.dhcpSetup()) {
    softSerial.println(F("DHCP failed"));
  } else {
    softSerial.println(F("DHCP successfull"));
  }

  softSerial.println(F("IP: "));
  printIp(ether.myip);
  softSerial.println(F("Netmask: "));
  printIp(ether.netmask);
  softSerial.println(F("Gateway: "));
  printIp(ether.gwip);
  softSerial.println(F("DNS: "));
  printIp(ether.dnsip);

}
