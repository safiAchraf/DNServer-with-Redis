
package main


import(
	"fmt"
	"time"
	"encoding/binary"
)
func ExtractTTL(response []byte) (time.Duration, error) {

	if len(response) < 12 {
		return 0, fmt.Errorf("response too short to extract TTL")
	}

	// Skip the header and question section to find the answer section
	offset := 12
	for {
		labelLen := int(response[offset])
		offset++
		if labelLen == 0 {
			break
		}
		offset += labelLen
	}
	offset += 4 // Skip QTYPE and QCLASS

	if len(response) < offset+4 {
		return 0, fmt.Errorf("response too short to extract TTL")
	}

	ttl := binary.BigEndian.Uint32(response[offset : offset+4])
	return time.Duration(ttl) * time.Second, nil
}