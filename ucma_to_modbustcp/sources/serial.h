SoftwareSerial softSerial(serialSoftwareRx, serialSoftwareTx);

void setupSerial() {
  softSerial.begin(9600);
}
