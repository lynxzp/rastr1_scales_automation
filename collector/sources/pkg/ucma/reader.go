package ucma

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

const (
	DataAccumAddr = 0x60
)

type Ucma struct {
	conn           net.Conn      `json:"-"`
	IP             string        `json:"ipaddr"`
	Port           string        `json:"-"`
	Rs485addr      uint8         `json:"rs485addr"`
	DataPerfAddr   uint16        `json:"dtype"`
	DataPerfValue  int32         `json:"DataPerfValue"`
	DataAccumValue int32         `json:"DataAccumValue"`
	Requests       int32         `json:"requests"`
	Responses      int32         `json:"responses"`
	RequestDelay   time.Duration `json:"-"`
	Id             int8          `json:"-"`
	Fraction       string        `json:"-"`
	time           time.Time     `json:"-"`
}

func (ucma *Ucma) Start(requestDelay time.Duration) {
	ucma.RequestDelay = requestDelay
	ucma.Port = "502"
	go ucma.read()
}

func (ucma *Ucma) read() {
	for {
		ucma.request()
		<-time.After(ucma.RequestDelay)
	}
}

func (ucma *Ucma) connect() (err error) {
	return err
}

func (ucma *Ucma) request() {
	firstRead := false
	ucma.requestProxy(ucma.DataPerfAddr, &ucma.DataPerfValue)
	if ucma.DataAccumValue == 0 {
		firstRead = true
	}
	ucma.requestProxy(DataAccumAddr, &ucma.DataAccumValue)
	if firstRead && (ucma.DataAccumValue != 0) {
		ucma.startSave()
	}
	ucma.periodicSave()
}

func (ucma *Ucma) requestProxy(addr uint16, dataP *int32) {
	if len(ucma.IP) == 0 {
		return
	}
	if ucma.Rs485addr == 0 {
		return
	}
	if ucma.DataPerfAddr == 0 {
		return
	}
	ucma.Requests++
	var err error
	ucma.conn, err = net.Dial("tcp", ucma.IP+":"+ucma.Port)
	if err != nil {
		log.Println("Can't connect to ", ucma.IP, ucma.Port)
		return
	}
	defer func() {
		ucma.conn.Close()
	}()

	data := ucma.modbusRequest(addr)
	if data >= 0 {
		ucma.Responses++
		*dataP = data
	}
}

func (ucma *Ucma) modbusRequest(addr uint16) int32 {
	// request
	foo := modbusRequest{501,
		0,
		6,
		ucma.Rs485addr,
		4,
		addr,
		2,
	}
	err := binary.Write(ucma.conn, binary.LittleEndian, foo)
	if err != nil {
		log.Fatal(err)
	}
	var data int32
	// response
	bytes := make([]byte, 64)
	for {
		n, conErr := ucma.conn.Read(bytes)
		if n != 0 {
			type response struct {
				Transaction uint32
				Unit        uint8
				Data        int32
			}
			var resp response
			jErr := json.Unmarshal(bytes[:n], &resp)
			if jErr != nil {
				log.Println("wrong scales response: ", jErr)
			}
			data = resp.Data
		}
		if conErr == io.EOF {
			break
		}
		if conErr != nil {
			log.Println("WW connection error:", conErr)
			return 0
		}
	}
	return data
}

type modbusRequest struct {
	transactionIdentifier uint16 // sequence
	protocolIdentifier    uint16 // = 0 always
	length                uint16 // = 6 for our message
	unitIdentifier        uint8  // slave address
	cmd                   uint8  // = 4 for our message
	dataAddress           uint16
	dataSize              uint16 // = 2, means 4 bytes
}
