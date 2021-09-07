class Modbus {
    struct modbusRequest{
        uint16_t transactionIdentifier;
        uint16_t protocolIdentifier;
        uint16_t length;
        uint8_t unitIdentifier;
        uint8_t cmd;
        uint16_t dataAddress;
        uint16_t dataSize;
    };

    struct modbusRequest* req;

    struct {
        uint16_t transactionIdentifier;
        uint16_t protocolIdentifier;
        uint16_t length;
        uint8_t unitIdentifier;
        uint8_t cmd;
        uint8_t  dataSize;
        int32_t  data;
    }modbusResponse;

public:
    uint8_t getRequestedDataAddress() {
        return req->dataAddress;
    }
    void setData(int32_t data) {
        modbusResponse.data = data;
    }
    char* getResponseBuf() {
        return (char*)(&modbusResponse);
    }
    uint8_t getResponseSize() {
        return sizeof(modbusResponse);
    }
    uint8_t getUnitIdentifier() {
        return modbusResponse.unitIdentifier;
    }
    bool encodeTCP(char *ptr) {
        softSerial.print(F("Incoming TCP:"));
        for(int i=0; i<12; i++){
          softSerial.print(ptr[i], HEX);
          softSerial.print(' ');
        }
        softSerial.println();

        req = (modbusRequest*)(ptr);
        //if(req->length != 7)
        //    return false;
        modbusResponse.transactionIdentifier = req->transactionIdentifier;
        modbusResponse.protocolIdentifier =    req->protocolIdentifier;
        modbusResponse.length =                req->length + 1;
        modbusResponse.unitIdentifier =        req->unitIdentifier;
        //    if((req->dataAddress != 0x60) &&
        //       (req->dataAddress != 0x5d) &&
        //       (req->dataAddress != 0x3f) &&
        //       (req->dataAddress != 0x44) &&
        //       (req->dataAddress != 0x37))
        //        return false;
        //    if(req->dataSize != 2)
        //        return false;
        modbusResponse.dataSize = 4;
        return true;
    }
};
