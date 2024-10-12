// pkg/tcp/tcp.go

package tcp

import (
    "encoding/binary"
    "errors"
    "net"
    "sync"

    "TCPIP_SUITE/pkg/utils"
)

// TCPHeader represents the TCP packet header.
type TCPHeader struct {
    SourcePort      uint16
    DestinationPort uint16
    SequenceNumber  uint32
    AcknowledgmentNumber uint32
    DataOffset      uint8 // 4 bits
    Reserved        uint8 // 3 bits
    Flags           uint8 // 9 bits (NS, CWR, ECE, URG, ACK, PSH, RST, SYN, FIN)
    WindowSize      uint16
    Checksum        uint16
    UrgentPointer   uint16
    Options         []byte
    Payload         []byte
}

// NewTCPHeader creates a new TCPHeader.
func NewTCPHeader(srcPort, destPort uint16, seqNum, ackNum uint32, flags uint8, window uint16, payload []byte) *TCPHeader {
    return &TCPHeader{
        SourcePort:      srcPort,
        DestinationPort: destPort,
        SequenceNumber:  seqNum,
        AcknowledgmentNumber: ackNum,
        DataOffset:      5, // No options
        Reserved:        0,
        Flags:           flags,
        WindowSize:      window,
        Checksum:        0, // To be calculated
        UrgentPointer:   0,
        Options:         []byte{},
        Payload:         payload,
    }
}

// Serialize converts the TCPHeader into bytes.
func (th *TCPHeader) Serialize(sourceIP, destIP net.IP) ([]byte, error) {
    headerLength := th.DataOffset * 4
    if headerLength < 20 {
        return nil, errors.New("invalid data offset, too short")
    }

    data := make([]byte, headerLength+len(th.Payload))
    binary.BigEndian.PutUint16(data[0:2], th.SourcePort)
    binary.BigEndian.PutUint16(data[2:4], th.DestinationPort)
    binary.BigEndian.PutUint32(data[4:8], th.SequenceNumber)
    binary.BigEndian.PutUint32(data[8:12], th.AcknowledgmentNumber)
    data[12] = (th.DataOffset << 4) | (th.Reserved << 1)
    data[13] = th.Flags
    binary.BigEndian.PutUint16(data[14:16], th.WindowSize)
    // Checksum to be calculated later
    binary.BigEndian.PutUint16(data[16:18], th.Checksum)
    binary.BigEndian.PutUint16(data[18:20], th.UrgentPointer)
    copy(data[20:headerLength], th.Options)
    copy(data[headerLength:], th.Payload)

    // Pseudo-header for checksum
    pseudoHeader := make([]byte, 12)
    copy(pseudoHeader[0:4], sourceIP.To4())
    copy(pseudoHeader[4:8], destIP.To4())
    pseudoHeader[8] = 0
    pseudoHeader[9] = 6 // TCP protocol number
    binary.BigEndian.PutUint16(pseudoHeader[10:12], uint16(len(data)))

    checksumData := append(pseudoHeader, data...)
    th.Checksum = CalculateChecksum(checksumData)
    binary.BigEndian.PutUint16(data[16:18], th.Checksum)

    return data, nil
}

// DeserializeTCPHeader converts bytes into a TCPHeader.
func DeserializeTCPHeader(data []byte) (*TCPHeader, error) {
    if len(data) < 20 {
        return nil, errors.New("data too short to be a valid TCP segment")
    }

    th := &TCPHeader{
        SourcePort:      binary.BigEndian.Uint16(data[0:2]),
        DestinationPort: binary.BigEndian.Uint16(data[2:4]),
        SequenceNumber:  binary.BigEndian.Uint32(data[4:8]),
        AcknowledgmentNumber: binary.BigEndian.Uint32(data[8:12]),
        DataOffset:      (data[12] >> 4) & 0x0F,
        Reserved:        (data[12] >> 1) & 0x07,
        Flags:           data[13],
        WindowSize:      binary.BigEndian.Uint16(data[14:16]),
        Checksum:        binary.BigEndian.Uint16(data[16:18]),
        UrgentPointer:   binary.BigEndian.Uint16(data[18:20]),
        Options:         []byte{},
        Payload:         []byte{},
    }

    headerLength := th.DataOffset * 4
    if len(data) < int(headerLength) {
        return nil, errors.New("data too short for header length")
    }

    if headerLength > 20 {
        th.Options = data[20:headerLength]
    }

    if len(data) > int(headerLength) {
        th.Payload = data[headerLength:]
    }

    return th, nil
}

// CalculateChecksum calculates the TCP checksum.
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

// TCPConnection represents a TCP connection state.
type TCPConnection struct {
    SourcePort      uint16
    DestinationPort uint16
    SourceIP        net.IP
    DestinationIP   net.IP
    SequenceNumber  uint32
    AcknowledgmentNumber uint32
    WindowSize      uint16
    State           string
    Mutex           sync.Mutex
}

// TCPManager manages TCP connections.
type TCPManager struct {
    Connections map[string]*TCPConnection
    Mutex       sync.Mutex
}

// NewTCPManager creates a new TCPManager.
func NewTCPManager() *TCPManager {
    return &TCPManager{
        Connections: make(map[string]*TCPConnection),
    }
}

// HandleSegment processes a TCP segment.
func (tm *TCPManager) HandleSegment(th *TCPHeader, sourceIP, destIP net.IP) {
    key := generateConnectionKey(sourceIP, destIP, th.SourcePort, th.DestinationPort)
    tm.Mutex.Lock()
    conn, exists := tm.Connections[key]
    if !exists {
        conn = &TCPConnection{
            SourcePort:      th.SourcePort,
            DestinationPort: th.DestinationPort,
            SourceIP:        sourceIP,
            DestinationIP:   destIP,
            SequenceNumber:  th.SequenceNumber,
            AcknowledgmentNumber: th.AcknowledgmentNumber,
            WindowSize:      th.WindowSize,
            State:           "CLOSED",
        }
        tm.Connections[key] = conn
    }
    tm.Mutex.Unlock()

    conn.Mutex.Lock()
    defer conn.Mutex.Unlock()

    // Simplified state machine
    switch conn.State {
    case "CLOSED":
        if th.Flags&0x02 != 0 { // SYN
            conn.State = "SYN_RECEIVED"
            conn.AcknowledgmentNumber = th.SequenceNumber + 1
            conn.SequenceNumber = 1001 // Example
            utils.Logger.Printf("Received SYN from %s:%d to %s:%d. Transitioning to SYN_RECEIVED.",
                sourceIP, th.SourcePort, destIP, th.DestinationPort)
            // Send SYN-ACK (not implemented)
        }
    case "SYN_RECEIVED":
        if th.Flags&0x10 != 0 { // ACK
            conn.State = "ESTABLISHED"
            utils.Logger.Printf("Received ACK from %s:%d to %s:%d. Transitioning to ESTABLISHED.",
                sourceIP, th.SourcePort, destIP, th.DestinationPort)
        }
    case "ESTABLISHED":
        if th.Flags&0x01 != 0 { // FIN
            conn.State = "CLOSE_WAIT"
            conn.AcknowledgmentNumber = th.SequenceNumber + 1
            utils.Logger.Printf("Received FIN from %s:%d to %s:%d. Transitioning to CLOSE_WAIT.",
                sourceIP, th.SourcePort, destIP, th.DestinationPort)
            // Send ACK and FIN (not implemented)
        }
    // Additional states can be handled here.
    default:
        utils.Logger.Printf("Unhandled TCP state: %s", conn.State)
    }
}

// generateConnectionKey generates a unique key for a TCP connection.
func generateConnectionKey(srcIP, destIP net.IP, srcPort, destPort uint16) string {
    return srcIP.String() + ":" + string(srcPort) + "->" + destIP.String() + ":" + string(destPort)
}