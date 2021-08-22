enum class data_t {
  performance=0x5d, // Производительность
  accumulation=0x60 // Накопление
};

class ucma {
public:
  
  static void setup() {
    Serial.begin(ucmaBaud,SERIAL_8N1);
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

  static int32_t read(uint8_t slave_addr, data_t datat) {
    uint8_t retries = ucmaRetries;
    int32_t resp = 0;
    while(retries--) {
      digitalWrite(ucmaDErePin,HIGH);
      request(slave_addr, datat);
      Serial.flush();
      digitalWrite(ucmaDErePin,LOW);
      resp = response(ucmaWaitResponseTimeoutMs);
      if (resp>0)
        return resp;
    }
    return resp;
  }

private:
  static unsigned char crc(unsigned char* buf, int cnt)
  {
  int i;
  unsigned char s=0;
  for (i = 1; i < cnt; i++) s+=buf[i];
  if (s>256) s=s-256;
  return s;
  }

  static void request(uint8_t slave_addr, data_t datat, uint8_t master_addr=1) {
    uint8_t buf[10];
    buf[0] = 0x54; // sync byte
    buf[1] = 0x08; // length
    buf[2] = slave_addr;
    buf[3] = master_addr;
    buf[4] = 0x01; // read cmd
    buf[5] = uint8_t(datat);
    buf[6] = 0;
    buf[7] = 0;
    buf[8] = 0;
    buf[9] = checksum(buf+1, 8);
    while(Serial.available()){
      softSerial.print(F("!read unexpected data:"));
      softSerial.println(int(Serial.read()));
    }
//    softSerial.print("0x");
//    for(uint8_t i=0; i<10; i++) {
//      char buffer[3]; sprintf (buffer, "%02x", buf[i]);
//      softSerial.print(buffer);
//    }
//    softSerial.print(" ");
    for(uint8_t i=0; i<10; i++) {
      Serial.print(char(buf[i]));
    }
  }

  static int32_t response(unsigned long timeout_ms) {
    unsigned long time = millis();
    uint8_t buf[10];
    uint8_t pos = 0;
    while(millis()-time < timeout_ms) {
      if(Serial.available()) {
        buf[pos] = Serial.read();
        if((pos>=6)&&(pos<=8)){
            char prbuf[20];
            sprintf(prbuf, "%02x ", int(buf[pos]));
            softSerial.print(prbuf);
        }
        pos++;
        if(pos>=10)
        {
          auto expc = checksum(buf+1,8);
          if (expc != buf[9]){
            softSerial.print("EE Wrong checksum. Expected:");
            softSerial.print(expc);
            softSerial.print("  received:");
            softSerial.println(buf[9]);
            return -1;
          }
          {
              char prbuf[10];
              uint16_t f = uint16_t(buf[6])*256+buf[7];
              char c=' ';
              if ((f>25000)&&(f<40000))
                  c='+';
              sprintf(prbuf, "% 6u%c ", f, c);
              softSerial.print(prbuf);
          }
          if( (buf[6]/16>=10) ||
          (buf[6]%16>=10) ||
          (buf[7]/16>=10) ||
          (buf[7]%16>=10) ||
          (buf[8]/16>=10) ||
          (buf[8]%16>=10))
              return -3;
          int32_t result = 0;
          result += buf[6]/16;
          result *= 10;
          result += buf[6]%16;
          result *= 10;
          result += buf[7]/16;
          result *= 10;
          result += buf[7]%16;
          result *= 10;
          result += buf[8]/16;
          result *= 10;
          result += buf[8]%16;
          return result;
        }
      }
    }
    return -2;
  }

  static uint8_t checksum(uint8_t* buf, int len)
  {
    uint8_t s=0;
    for (uint8_t i = 0; i < len; i++) 
      s += buf[i];
    return s;
  }

};
