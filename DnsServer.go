package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type DNSHeader struct {
    ID      uint16 
    Flags   uint16 
    QDCount uint16 
    ANCount uint16 
    NSCount uint16 
    ARCount uint16 
}

func (h *DNSHeader) ToBytes() ([]byte, error) {
    buffer := make([]byte, 12)
    binary.BigEndian.PutUint16(buffer[0:2], h.ID)
    binary.BigEndian.PutUint16(buffer[2:4], h.Flags)
    binary.BigEndian.PutUint16(buffer[4:6], h.QDCount)
    binary.BigEndian.PutUint16(buffer[6:8], h.ANCount)
    binary.BigEndian.PutUint16(buffer[8:10], h.NSCount)
    binary.BigEndian.PutUint16(buffer[10:12], h.ARCount)
    return buffer, nil
}

func (h *DNSHeader) DnsHeaderFromBytes(data []byte) error {
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

type DNSquestion struct {
    domain string
    QuesType uint16
    QuesClass uint16
}


func (q *DNSquestion) ToBytes() ([]byte , error){
    DomainParts := []byte{}

    for _,part := range strings.Split(q.domain, ".") {
        DomainParts = append(DomainParts, byte(len(part)))
        DomainParts = append(DomainParts, []byte(part)...)
    }

    DomainParts = append(DomainParts, byte(0))

    buffer := make([]byte , len(DomainParts) + 4 )

    copy(buffer, DomainParts)

    binary.BigEndian.PutUint16(buffer[len(DomainParts):len(DomainParts)+2], q.QuesType)
    binary.BigEndian.PutUint16(buffer[len(DomainParts)+2:], q.QuesClass)
    return buffer, nil
}

func (q * DNSquestion) DnsQuestionFromBytes(rawBytes [] byte) error {
    q.domain = ""
    i := 0

    for i < len(rawBytes) && rawBytes[i] != 0 {
        partLength := int(rawBytes[i])
        i++

        if i + partLength > len(rawBytes) {
            return fmt.Errorf("the data u sent is not complete (the domain size dont match with the content)")
        }

        if len(q.domain) > 0 {
            q.domain += "."
        }

        q.domain = string(rawBytes[i : i+ partLength])
        i += partLength
    }
    i++

    if len(rawBytes) < i+4 {
        return fmt.Errorf("the class and type field are not recieved or corrupted")
    }

    q.QuesType = binary.BigEndian.Uint16(rawBytes[i : i+2])

    i+=2
    q.QuesClass = binary.BigEndian.Uint16(rawBytes[i : i+2])

    return nil
}


type DNSAnswer struct {
    Name       string 
    Type       uint16 
    Class      uint16 
    TTL        uint32 
    DataLength uint16 
    RData      []byte 
}

func (a *DNSAnswer) ToBytes() ([]byte, error) {
    nameParts := []byte{}
    for _, part := range strings.Split(a.Name, ".") {
        nameParts = append(nameParts, byte(len(part)))
        nameParts = append(nameParts, []byte(part)...)
    }
    nameParts = append(nameParts, 0) 

    buffer := make([]byte, len(nameParts)+10+int(a.DataLength))
    copy(buffer, nameParts)

    offset := len(nameParts)
    binary.BigEndian.PutUint16(buffer[offset:offset+2], a.Type)
    offset += 2
    binary.BigEndian.PutUint16(buffer[offset:offset+2], a.Class)
    offset += 2
    binary.BigEndian.PutUint32(buffer[offset:offset+4], a.TTL)
    offset += 4
    binary.BigEndian.PutUint16(buffer[offset:offset+2], a.DataLength)
    offset += 2
    copy(buffer[offset:], a.RData)

    return buffer, nil
}


func (a *DNSAnswer) DnsAnswerFromBytes(data []byte) error {
    i := 0
    a.Name = ""
    
    for i < len(data) && data[i] != 0 {
        labelLength := int(data[i])
        i++
        if i+labelLength > len(data) {
            return fmt.Errorf("invalid DNS answer format: label length exceeds data size")
        }
        if len(a.Name) > 0 {
            a.Name += "."
        }
        a.Name += string(data[i : i+labelLength])
        i += labelLength
    }

    i++ 

    if len(data) < i+10 {
        return fmt.Errorf("invalid DNS answer format: insufficient data")
    }

    a.Type = binary.BigEndian.Uint16(data[i : i+2])
    i += 2
    a.Class = binary.BigEndian.Uint16(data[i : i+2])
    i += 2
    a.TTL = binary.BigEndian.Uint32(data[i : i+4])
    i += 4
    a.DataLength = binary.BigEndian.Uint16(data[i : i+2])
    i += 2

    if len(data) < i+int(a.DataLength) {
        return fmt.Errorf("invalid DNS answer format: insufficient data for RData")
    }

    a.RData = data[i : i+int(a.DataLength)]

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

		message , err := header.ToBytes()

		n , err = conn.WriteToUDP(message , senderAddr)
		if err != nil {
			fmt.Println("hello negro didnt sent")
		}
	}
}