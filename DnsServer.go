package main

import (
	"fmt"
	"net"
	"encoding/binary"
)


type DNSBody struct {

}
type DNSHeader struct {
    ID      uint16 
    Flags   uint16 
    QDCount uint16 
    ANCount uint16 
    NSCount uint16 
    ARCount uint16 
}

func (h *DNSHeader) DnsHeaderToBytes() ([]byte, error) {
    buffer := make([]byte, 12)
    binary.BigEndian.PutUint16(buffer[0:2], h.ID)
    binary.BigEndian.PutUint16(buffer[2:4], h.Flags)
    binary.BigEndian.PutUint16(buffer[4:6], h.QDCount)
    binary.BigEndian.PutUint16(buffer[6:8], h.ANCount)
    binary.BigEndian.PutUint16(buffer[8:10], h.NSCount)
    binary.BigEndian.PutUint16(buffer[10:12], h.ARCount)
    return buffer, nil
}

func (h *DNSHeader) SerializeDnsHeaderFromBytes(data []byte) error {
    if len(data) < 12 {
        return fmt.Errorf("data too short to be a valid DNS header")
    }
    h.ID = binary.BigEndian.Uint16(data[0:2])
    h.Flags = binary.BigEndian.Uint16(data[2:4])
    h.QDCount = binary.BigEndian.Uint16(data[4:6])
    h.ANCount = binary.BigEndian.Uint16(data[6:8])
    h.NSCount = binary.BigEndian.Uint16(data[8:10])
    h.ARCount = binary.BigEndian.Uint16(data[10:12])
    return nil
}

func main() {
	addr , err := net.ResolveUDPAddr("udp" , ":8080")
	if err != nil {
		fmt.Println("error creating the addr")
	}
	conn , err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Println("error listening")
	}

	buffer := make([]byte , 1024)
	header := DNSHeader{
        ID:      0x1234,
        Flags:   0x8180, 
        QDCount: 1,      
        ANCount: 1,      
        NSCount: 0,
        ARCount: 0,
    }

	for {
		n , senderAddr , err := conn.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("error recieve packet")
			continue
		}
		fmt.Println("recieved msg from %s : %s 	\n" , senderAddr , string(buffer[:n]))

		message , err := header.DnsHeaderToBytes()

		n , err = conn.WriteToUDP(message , senderAddr)
		if err != nil {
			fmt.Println("hello negreo didnt sent")
		}
	}
}