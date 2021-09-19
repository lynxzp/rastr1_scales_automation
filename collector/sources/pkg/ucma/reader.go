package ucma

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const ScalsesNums = 15

type Ucma struct {
	conn         net.Conn		`json:"-"`
	IP           string			`json:"-"`
	Port         string			`json:"-"`
	Rs485addr    uint8			`json:"-"`
	DType        uint16			`json:"-"`
	Data		 int32		    `json:"data"`
	Ready        bool   		`json:"ready"`
	RequestDelay time.Duration	`json:"-"`
	Mu 			 sync.Mutex		`json:"-"`
}

func (ucma *Ucma) Start(requestDelay time.Duration) {
	ucma.Mu.Lock()
	ucma.RequestDelay = requestDelay
	ucma.Port = "502"
	ucma.Mu.Unlock()
	go ucma.read()
}

func (ucma *Ucma) read() {
	for {
		ucma.Mu.Lock()
		if ucma.Ready == true {
			ucma.Mu.Unlock()
			ucma.request()
		} else {
		ucma.Mu.Unlock()
		}
		<-time.After(ucma.RequestDelay)
	}
}

func (ucma *Ucma) connect() (err error) {
	return err
}

func (ucma *Ucma) request() {
	var err error
	ucma.Mu.Lock()
	ucma.conn, err = net.Dial("tcp", ucma.IP+":"+ucma.Port)
	ucma.Mu.Unlock()
	if err!=nil {
		log.Println("Can't connect to ", ucma.IP, ucma.Port)
		return
	}
	// request
	ucma.Mu.Lock()
	foo := modbusRequest{501,
		0,
		6,
		ucma.Rs485addr,
		4,
		ucma.DType,
		2,
	}
	err = binary.Write(ucma.conn, binary.LittleEndian, foo)
	ucma.Mu.Unlock()
	if err != nil {
		log.Fatal(err)
	}

	// response
	bytes := make([]byte, 64)
	for {
		ucma.Mu.Lock()
		n, conErr := ucma.conn.Read(bytes)
		ucma.Mu.Unlock()
		if n != 0 {
			type response struct {
				Transaction uint32
				Unit		uint8
				Data		int32
			}
			var resp response
			jErr := json.Unmarshal(bytes[:n], &resp)
			if jErr!= nil {log.Println("wrong scales response: ", jErr)}
			ucma.Mu.Lock()
			ucma.Data = resp.Data
			ucma.Mu.Unlock()
		}
		if conErr == io.EOF {
			break
		}
	}
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
