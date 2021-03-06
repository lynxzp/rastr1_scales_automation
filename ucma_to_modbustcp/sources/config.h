constexpr uint8_t ethernetMAC[] = { 0x74,0x69,0x69,0x2D,0x30,0x1F };
constexpr uint8_t ethernetCS = 10;
constexpr uint16_t ethernetBufferSize = 390; // 340 is minimum for DHCP + 50 just for case
constexpr uint16_t modbusTcpPort = 502;

constexpr uint16_t serialBaud = 9600;
constexpr uint8_t serialSoftwareRx = A2;
constexpr uint8_t serialSoftwareTx = A3;

constexpr uint16_t ucmaBaud = 9600;
constexpr uint32_t ucmaWaitResponseTimeoutMs = 500;
constexpr uint8_t ucmaRetries = 1;
constexpr uint8_t ucmaDErePin = 2;
