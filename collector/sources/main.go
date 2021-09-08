package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)
import "fmt"
import "bufio"
import "os"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	startClient("192.168.1.12:502")
	startClient("127.0.0.1:5030")
}

func startClient(address string) {
	//connect to this socket
	connClient, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
	}

	go read(connClient)
	go periodicRequest(connClient)

	for {

		//read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}

		//send to socket
		fmt.Fprint(connClient, text)

	}
}

func read(reader io.Reader) {
	bytes := make([]byte, 64)
	for {
		n, err := reader.Read(bytes)
		fmt.Printf("%v\n%q\n", bytes[:n],bytes[:n])
		if err == io.EOF {
			break
		}
	}
}

func periodicRequest(writer io.Writer) {
	for {
		go request(writer)
		<-time.After(1000 * time.Second)
	}
}

type modbusRequest struct {
	transactionIdentifier uint16	// sequence
	protocolIdentifier    uint16	// always = 0
	length                uint16	// = 6 for our message
	unitIdentifier        uint8		// slave address
	cmd                   uint8		// = 4 for our message
	dataAddress           uint16
	dataSize              uint16	// = 2, means 4 bytes
}

func request(writer io.Writer) {
	foo := modbusRequest{501,
		0,
		6,
		2,
		4,
		0x60,
		2,
	}
	err := binary.Write(writer, binary.LittleEndian, foo)
	if err != nil {
		log.Fatal(err)
	}
	//bytes := make([]byte, 64)
	//_, err = writer.Write(bytes)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
