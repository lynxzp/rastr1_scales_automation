package ucma

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Ucma struct {
	conn net.Conn
	ip string
	port string
	rs485addr uint8
	data         int
	requestDelay time.Duration
}

func (ucma *Ucma)Start(ip string, port string, requestDelay time.Duration) {
	ucma.ip = ip
	ucma.port = port
	ucma.requestDelay = requestDelay
	go ucma.read()
}

func (ucma *Ucma)read() {
	for {
		ucma.request()
		<-time.After(ucma.requestDelay)
	}
}

func (ucma *Ucma)connect() (err error) {
	return err
}

func (ucma *Ucma)request() {
	ucma.conn, _ = net.Dial("tcp", ucma.ip+":"+ucma.port)
	// request
	foo := modbusRequest{501,
		0,
		6,
		2,
		4,
		0x60,
		2,
	}
	err := binary.Write(ucma.conn, binary.LittleEndian, foo)
	if err != nil {
		log.Fatal(err)
	}

	// response
	bytes := make([]byte, 64)
	for {
		n, err := ucma.conn.Read(bytes)
		if n != 0 {
			fmt.Printf("%v\n%q\n", bytes[:n],bytes[:n])
			// todo: decode and write to data variable
		}
		if err == io.EOF {
			break
		}
	}
}

type modbusRequest struct {
	transactionIdentifier uint16	// sequence
	protocolIdentifier    uint16	// = 0 always
	length                uint16	// = 6 for our message
	unitIdentifier        uint8		// slave address
	cmd                   uint8		// = 4 for our message
	dataAddress           uint16
	dataSize              uint16	// = 2, means 4 bytes
}

