package main

import (
	"encoding/binary"
	"fmt"
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

type DNSQuery struct {
	Header          DNSHeader
	Questions       []DNSQuestion
}

type DNSQuestion struct {
    domain string
    QuesType uint16
    QuesClass uint16
}

type DNSResourceRecords  struct {
    Name       string 
    Type       uint16 
    Class      uint16 
    TTL        uint32 
    DataLength uint16 
    RData      []byte 
}


type DNSMessage struct {
	Header          DNSHeader
	Questions       []DNSQuestion
	ResourceRecords []DNSResourceRecords
}



func (answer *DNSResourceRecords ) serialize() ([]byte, error) {
    nameParts := []byte{}
    for _, part := range strings.Split(answer.Name, ".") {
        nameParts = append(nameParts, byte(len(part)))
        nameParts = append(nameParts, []byte(part)...)
    }
    nameParts = append(nameParts, 0) 

    buffer := make([]byte, len(nameParts)+10+int(answer.DataLength))
    copy(buffer, nameParts)

    offset := len(nameParts)
    binary.BigEndian.PutUint16(buffer[offset:offset+2], answer.Type)
    offset += 2
    binary.BigEndian.PutUint16(buffer[offset:offset+2], answer.Class)
    offset += 2
    binary.BigEndian.PutUint32(buffer[offset:offset+4], answer.TTL)
    offset += 4
    binary.BigEndian.PutUint16(buffer[offset:offset+2], answer.DataLength)
    offset += 2
    copy(buffer[offset:], answer.RData)

    return buffer, nil
}


func (answer *DNSResourceRecords) FromBytes(data []byte) error {
    i := 0
    answer.Name = ""
    
    for i < len(data) && data[i] != 0 {
        labelLength := int(data[i])
        i++
        if i+labelLength > len(data) {
            return fmt.Errorf("invalid DNS answer format: label length exceeds data size")
        }
        if len(answer.Name) > 0 {
            answer.Name += "."
        }
        answer.Name += string(data[i : i+labelLength])
        i += labelLength
    }

    i++ 

    if len(data) < i+10 {
        return fmt.Errorf("invalid DNS answer format: insufficient data")
    }

    answer.Type = binary.BigEndian.Uint16(data[i : i+2])
    i += 2
    answer.Class = binary.BigEndian.Uint16(data[i : i+2])
    i += 2
    answer.TTL = binary.BigEndian.Uint32(data[i : i+4])
    i += 4
    answer.DataLength = binary.BigEndian.Uint16(data[i : i+2])
    i += 2

    if len(data) < i+int(answer.DataLength) {
        return fmt.Errorf("invalid DNS answer format: insufficient data for RData")
    }

    answer.RData = data[i : i+int(answer.DataLength)]

    return nil
}




func (q *DNSQuestion) serialize() ([]byte , error){
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

func (q * DNSQuestion) DnsQuestionFromBytes(rawBytes [] byte) error {
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

        q.domain += string(rawBytes[i : i+ partLength])
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




func (query *DNSQuery) DNSMessageFromBytes(data []byte) error {
    var offset int
    var err error

    if err = query.Header.DnsHeaderFromBytes(data[:12]); err != nil {
        return fmt.Errorf("failed to parse DNS header: %v", err)
    }
    offset = 12

    query.Questions = make([]DNSQuestion, query.Header.QDCount)

    for i := 0; i < int(query.Header.QDCount); i++ {

        var question DNSQuestion

        if err = question.DnsQuestionFromBytes(data[offset:]); err != nil {
            return fmt.Errorf("failed to parse DNS question: %v", err)
        }
        query.Questions[i] = question

        offset += len(question.domain) + 2 + 2 + 1 
    }

    return nil
}




func (h *DNSHeader) serialize() ([]byte, error) {
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



func (message *DNSMessage) serialize() ([]byte , error){
    buffer := []byte{}

    header , err := message.Header.serialize()
    if err != nil {
        return []byte{} , fmt.Errorf("Error in serializing the header")
    }
    buffer = append(buffer, header... )

    for _ , question := range message.Questions {
        questionSerialized , err := question.serialize()
        if err != nil {
            return []byte{} , fmt.Errorf("Error in serializing the Question with the Domain : %v " , question.domain) 
        }
        buffer = append(buffer, questionSerialized...)
    }

    for _ , answer := range message.ResourceRecords {
        RecourceRecords  , err := answer.serialize()
        if err != nil {
            return []byte{} , fmt.Errorf("Error in serializing the Answer with the Name : %v " , answer.Name)
        }
        buffer = append(buffer, RecourceRecords...)
    }

    return buffer , nil
}

func (message *DNSMessage) DNSMessageFromBytes(data []byte) error {
    var offset int
    var err error

    if err = message.Header.DnsHeaderFromBytes(data[:12]); err != nil {
        return fmt.Errorf("failed to parse DNS header: %v", err)
    }
    offset = 12

    message.Questions = make([]DNSQuestion, int(message.Header.QDCount))

    for i := 0; i < int(message.Header.QDCount); i++ {

        var question DNSQuestion

        if err = question.DnsQuestionFromBytes(data[offset:]); err != nil {
            fmt.Printf("failed to parse DNS question: %v", err)
            return fmt.Errorf("failed to parse DNS question: %v", err)
        }
        message.Questions[i] = question

        offset += len(question.domain) + 2 + 2 + 1 
    }

    message.ResourceRecords = make([]DNSResourceRecords, message.Header.ANCount)

    for i := 0; i < int(message.Header.ANCount); i++ {

        var resourceRecord DNSResourceRecords

        if err = resourceRecord.FromBytes(data[offset:]); err != nil {
            fmt.Printf("failed to parse DNS resource record: %v", err)
            return fmt.Errorf("failed to parse DNS resource record: %v", err)
        }
        message.ResourceRecords[i] = resourceRecord
        offset += len(resourceRecord.Name) + 2 + 2 + 4 + 2 + int(resourceRecord.DataLength) // Name + TYPE + CLASS + TTL + DataLength + RData
    }

    return nil
}