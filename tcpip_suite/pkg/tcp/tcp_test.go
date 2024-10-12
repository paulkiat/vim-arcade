// pkg/tcp/tcp_test.go

package tcp

import (
    "net"
    "testing"
)

func TestTCPHeaderSerialization(t *testing.T) {
    sourceIP := net.ParseIP("221.54.8.1")
    destIP := net.ParseIP("000.54.0.6")
    protocol := uint8(6) // TCP
    payload := []byte("Test TCP Payload")

    ih := NewIPv4Header(sourceIP, destIP, protocol, payload)
    serializedIP, err := ih.Serialize()
    if err != nil {
        t.Fatalf("Failed to serialize IPv4 header: %v", err)
    }

    tcpHeader := NewTCPHeader(12345, 80, 1000, 0, 0x02, 65535, []byte{})
    serializedTCP, err := tcpHeader.Serialize(sourceIP, destIP)
    if err != nil {
        t.Fatalf("Failed to serialize TCP header: %v", err)
    }

    deserializedTCP, err := DeserializeTCPHeader(serializedTCP)
    if err != nil {
        t.Fatalf("Failed to deserialize TCP header: %v", err)
    }

    if deserializedTCP.SourcePort != tcpHeader.SourcePort {
        t.Errorf("SourcePort mismatch. Expected %d, got %d", tcpHeader.SourcePort, deserializedTCP.SourcePort)
    }

    if deserializedTCP.DestinationPort != tcpHeader.DestinationPort {
        t.Errorf("DestinationPort mismatch. Expected %d, got %d", tcpHeader.DestinationPort, deserializedTCP.DestinationPort)
    }

    if deserializedTCP.SequenceNumber != tcpHeader.SequenceNumber {
        t.Errorf("SequenceNumber mismatch. Expected %d, got %d", tcpHeader.SequenceNumber, deserializedTCP.SequenceNumber)
    }

    if deserializedTCP.AcknowledgmentNumber != tcpHeader.AcknowledgmentNumber {
        t.Errorf("AcknowledgmentNumber mismatch. Expected %d, got %d", tcpHeader.AcknowledgmentNumber, deserializedTCP.AcknowledgmentNumber)
    }

    if deserializedTCP.Flags != tcpHeader.Flags {
        t.Errorf("Flags mismatch. Expected 0x%X, got 0x%X", tcpHeader.Flags, deserializedTCP.Flags)
    }

    if deserializedTCP.WindowSize != tcpHeader.WindowSize {
        t.Errorf("WindowSize mismatch. Expected %d, got %d", tcpHeader.WindowSize, deserializedTCP.WindowSize)
    }

    if len(deserializedTCP.Options) != len(tcpHeader.Options) {
        t.Errorf("Options length mismatch. Expected %d, got %d", len(tcpHeader.Options), len(deserializedTCP.Options))
    }
}
