// pkg/udp/udp.go

package udp

import (
    "encoding/binary"
    "errors"
    "net"

    "TCPIP_SUITE/pkg/utils"
)

// UDPHeader represents the UDP packet header.
type UDPHeader struct {
    SourcePort      uint16
    DestinationPort uint16
    Length          uint16
    Checksum        uint16
    Payload         []byte
}

// NewUDPHeader creates a new UDPHeader.
func NewUDPHeader(srcPort, destPort uint16, payload []byte) *UDPHeader {
    length := uint16(8 + len(payload)) // 8 bytes header + payload
    return &UDPHeader{
        SourcePort:      srcPort,
        DestinationPort: destPort,
        Length:          length,
        Checksum:        0, // To be calculated
        Payload:         payload,
    }
}

// Serialize converts the UDPHeader into bytes.
func (uh *UDPHeader) Serialize(sourceIP, destIP net.IP) ([]byte, error) {
    if sourceIP == nil || destIP == nil {
        return nil, errors.New("invalid IP addresses")
    }

    data := make([]byte, uh.Length)
    binary.BigEndian.PutUint16(data[0:2], uh.SourcePort)
    binary.BigEndian.PutUint16(data[2:4], uh.DestinationPort)
    binary.BigEndian.PutUint16(data[4:6], uh.Length)
    // Checksum to be calculated later
    binary.BigEndian.PutUint16(data[6:8], uh.Checksum)
    copy(data[8:], uh.Payload)

    // Pseudo-header for checksum
    pseudoHeader := make([]byte, 12)
    copy(pseudoHeader[0:4], sourceIP.To4())
    copy(pseudoHeader[4:8], destIP.To4())
    pseudoHeader[8] = 0
    pseudoHeader[9] = 17 // UDP protocol number
    binary.BigEndian.PutUint16(pseudoHeader[10:12], uh.Length)

    checksumData := append(pseudoHeader, data...)
    uh.Checksum = CalculateChecksum(checksumData)
    binary.BigEndian.PutUint16(data[6:8], uh.Checksum)

    return data, nil
}

// DeserializeUDPHeader converts bytes into a UDPHeader.
func DeserializeUDPHeader(data []byte) (*UDPHeader, error) {
    if len(data) < 8 {
        return nil, errors.New("data too short to be a valid UDP segment")
    }

    uh := &UDPHeader{
        SourcePort:      binary.BigEndian.Uint16(data[0:2]),
        DestinationPort: binary.BigEndian.Uint16(data[2:4]),
        Length:          binary.BigEndian.Uint16(data[4:6]),
        Checksum:        binary.BigEndian.Uint16(data[6:8]),
        Payload:         data[8:],
    }

    // Verify checksum (optional)
    // Implement checksum verification if needed.

    return uh, nil
}

// CalculateChecksum calculates the UDP checksum.
func CalculateChecksum(data []byte) uint16 {
    var sum uint32
    for i := 0; i+1 < len(data); i += 2 {
        sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
    }
    if len(data)%2 != 0 {
        sum += uint32(data[len(data)-1]) << 8
    }
    for sum > 0xFFFF {
        sum = (sum >> 16) + (sum & 0xFFFF)
    }
    return ^uint16(sum)
}