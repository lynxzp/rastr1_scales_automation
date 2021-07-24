
byte Ethernet::buffer[ethernetBufferSize];

void setupEthernet() {
  Serial.begin(9600);
  Serial.println(F("\n  --== testDHCP ==--"));

  Serial.print("MAC: ");
  for (byte i = 0; i < 6; ++i) {
    Serial.print(ethernetMAC[i], HEX);
    if (i < 5)
      Serial.print(':');
  }
  Serial.println();

  if (ether.begin(sizeof Ethernet::buffer, ethernetMAC, ethernetCS) == 0){
    Serial.println(F("Failed to access Ethernet controller"));
  } else {
    Serial.println(F("Sucessfull access Ethernet controller"));
  }

  Serial.println(F("Setting up DHCP"));
  if (!ether.dhcpSetup()) {
    Serial.println(F("DHCP failed"));
  } else {
    Serial.println(F("DHCP FINISHED"));
  }


  ether.printIp("My IP: ", ether.myip);
  ether.printIp("Netmask: ", ether.netmask);
  ether.printIp("GW IP: ", ether.gwip);
  ether.printIp("DNS IP: ", ether.dnsip);

}
