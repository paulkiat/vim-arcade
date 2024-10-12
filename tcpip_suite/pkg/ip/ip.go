// pkg/ip/ip.go

package ip

import (
    "encoding/binary"
    "errors"
    "net"
    "sync"

    "TCPIP_SUITE/pkg/utils"
)

// IPv4Header represents the IPv4 packet header.
type IPv4Header struct {
    Version        uint8
    IHL            uint8
    TypeOfService  uint8
    TotalLength    uint16
    Identification uint16
    Flags          uint8
    FragmentOffset uint16
    TTL            uint8
    Protocol       uint8
    HeaderChecksum uint16
    SourceIP       net.IP
    DestinationIP  net.IP
    Options        []byte
    Payload        []byte
}

// NewIPv4Header creates a new IPv4Header.
func NewIPv4Header(sourceIP, destIP net.IP, protocol uint8, payload []byte) *IPv4Header {
    ihl := 5 // No options
    totalLength := uint16(20 + len(payload)) // 20 bytes header + payload

    return &IPv4Header{
        Version:        4,
        IHL:            ihl,
        TypeOfService:  0,
        TotalLength:    totalLength,
        Identification: 0,
        Flags:          0,
        FragmentOffset: 0,
        TTL:            64,
        Protocol:       protocol,
        HeaderChecksum: 0, // To be calculated
        SourceIP:       sourceIP.To4(),
        DestinationIP:  destIP.To4(),
        Options:        []byte{},
        Payload:        payload,
    }
}

// Serialize converts the IPv4Header into bytes.
func (ih *IPv4Header) Serialize() ([]byte, error) {
    if ih.SourceIP == nil || ih.DestinationIP == nil {
        return nil, errors.New("invalid IP addresses")
    }

    headerLength := ih.IHL * 4
    if headerLength < 20 {
        return nil, errors.New("invalid IHL, too short")
    }

    totalLength := ih.TotalLength
    if int(totalLength) < headerLength {
        return nil, errors.New("total length less than header length")
    }

    data := make([]byte, totalLength)
    data[0] = (ih.Version << 4) | ih.IHL
    data[1] = ih.TypeOfService
    binary.BigEndian.PutUint16(data[2:4], ih.TotalLength)
    binary.BigEndian.PutUint16(data[4:6], ih.Identification)
    flagsAndFragment := (uint16(ih.Flags) << 13) | ih.FragmentOffset
    binary.BigEndian.PutUint16(data[6:8], flagsAndFragment)
    data[8] = ih.TTL
    data[9] = ih.Protocol
    // Header checksum to be calculated later
    binary.BigEndian.PutUint16(data[10:12], ih.HeaderChecksum)
    copy(data[12:16], ih.SourceIP)
    copy(data[16:20], ih.DestinationIP)
    copy(data[20:headerLength], ih.Options)
    copy(data[headerLength:], ih.Payload)

    // Calculate checksum
    ih.HeaderChecksum = CalculateChecksum(data[:headerLength])
    binary.BigEndian.PutUint16(data[10:12], ih.HeaderChecksum)

    return data, nil
}

// DeserializeIPv4Header converts bytes into an IPv4Header.
func DeserializeIPv4Header(data []byte) (*IPv4Header, error) {
    if len(data) < 20 {
        return nil, errors.New("data too short to be a valid IPv4 packet")
    }

    version := data[0] >> 4
    ihl := data[0] & 0x0F
    if version != 4 {
        return nil, errors.New("not an IPv4 packet")
    }
    headerLength := ihl * 4
    if len(data) < int(headerLength) {
        return nil, errors.New("data too short for header length")
    }

    totalLength := binary.BigEndian.Uint16(data[2:4])
    if len(data) < int(totalLength) {
        return nil, errors.New("data too short for total length")
    }

    ih := &IPv4Header{
        Version:        version,
        IHL:            ihl,
        TypeOfService:  data[1],
        TotalLength:    totalLength,
        Identification: binary.BigEndian.Uint16(data[4:6]),
        Flags:          uint8(data[6] >> 5),
        FragmentOffset: binary.BigEndian.Uint16(data[6:8]) & 0x1FFF,
        TTL:            data[8],
        Protocol:       data[9],
        HeaderChecksum: binary.BigEndian.Uint16(data[10:12]),
        SourceIP:       net.IP(data[12:16]),
        DestinationIP:  net.IP(data[16:20]),
        Options:        data[20:headerLength],
        Payload:        data[headerLength:totalLength],
    }

    // Verify checksum
    calculatedChecksum := CalculateChecksum(data[:headerLength])
    if calculatedChecksum != ih.HeaderChecksum {
        utils.Logger.Printf("Checksum mismatch: calculated 0x%X, expected 0x%X", calculatedChecksum, ih.HeaderChecksum)
        return nil, errors.New("checksum mismatch")
    }

    return ih, nil
}

// CalculateChecksum calculates the IPv4 header checksum.
func CalculateChecksum(data []byte) uint16 {
    var sum uint32
    for i := 0; i < len(data); i += 2 {
        word := binary.BigEndian.Uint16(data[i : i+2])
        sum += uint32(word)
    }
    for sum > 0xFFFF {
        sum = (sum >> 16) + (sum & 0xFFFF)
    }
    return ^uint16(sum)
}

// HandleFragmentation handles IP fragmentation if needed.
// Placeholder implementation.
func HandleFragmentation(ih *IPv4Header) ([]*IPv4Header, error) {
    // Implement fragmentation logic based on MTU and payload size.
    return []*IPv4Header{ih}, nil
}