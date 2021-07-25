enum class data_t {
  performance, // Производительность
  accumulation // Накопление
};

class ucma {
public:
  
  static void setup() {
    Serial.begin(ucmaBaud);
    pinMode(ucmaDErePin,OUTPUT);
  }
  
  static bool available() {
    return Serial.available();
  }
  
  static uint8_t read() {
    return Serial.read();
  }
  
  static void write(uint8_t c) {
    Serial.write(c);
  }

  static uint8_t read(uint8_t slave_addr, uint8_t master_addr, data_t datat) {
    uint8_t retries = ucmaRetries;
    while(retries--) {
      digitalWrite(ucmaDErePin,HIGH);
      request(slave_addr, master_addr, datat);
      request(slave_addr, master_addr, datat);
      digitalWrite(ucmaDErePin,LOW);
      auto resp = response(ucmaWaitResponseTimeoutMs);
      if (resp>0)
        return resp;
    }
    return -1;
  }

private:

  static void request(uint8_t slave_addr, uint8_t master_addr, data_t datat) {
    uint8_t buf[10];
    buf[0] = 0x54; // sync byte
    buf[1] = 0x08; // length
    buf[2] = slave_addr;
    buf[3] = master_addr;
    buf[4] = 0x01; // read cmd
    
    buf[5] = 0x5d; // data_t::performance
    if (datat==data_t::accumulation)
        buf[5] = 0x60;
        
    buf[6] = 0;
    buf[7] = 0;
    buf[8] = 0;
    buf[9] = checksum(buf+1, 8);
    while(Serial.available()){
      softSerial.print(F("!read unexpected data:"));
      softSerial.println(int(Serial.read()));
    }
    for(uint8_t i=0; i<10; i++) {
      Serial.print((char)buf[i]);
    }
  }

  static int16_t response(unsigned long timeout_ms) {
    unsigned long time = millis();
    uint8_t buf[10];
    uint8_t pos = 0;
    while(millis()-time < timeout_ms) {
      if(Serial.available()) {
        buf[pos] = Serial.read();
        pos++;
        if(pos>=10)
        {
          if (checksum(buf+1,8) != buf[9])
            return -1;
          return (buf[6]-'0')*100 + (buf[7]-'0')*10 + buf[8]-'0';
        }
        softSerial.print("incoming: ");
        softSerial.println(int(buf[pos-1]));
      }
    }
    return -1;
  }

  static uint8_t checksum(uint8_t* buf, int len)
  {
    uint8_t s=0;
    for (uint8_t i = 0; i < len; i++) 
      s += buf[i];
    return s;
  }

};
