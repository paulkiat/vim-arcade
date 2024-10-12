// pkg/udp/udp_test.go

package udp

import (
    "net"
    "testing"
)

func TestUDPHeaderSerialization(t *testing.T) {
    sourceIP := net.ParseIP("192.168.1.2")
    destIP := net.ParseIP("192.168.1.3")
    srcPort := uint16(12345)
    destPort := uint16(80)
    payload := []byte("Test UDP Payload")

    udpHeader := NewUDPHeader(srcPort, destPort, payload)
    serialized, err := udpHeader.Serialize(sourceIP, destIP)
    if err != nil {
        t.Fatalf("Failed to serialize UDP header: %v", err)
    }

    deserialized, err := DeserializeUDPHeader(serialized)
    if err != nil {
        t.Fatalf("Failed to deserialize UDP header: %v", err)
    }

    if deserialized.SourcePort != udpHeader.SourcePort {
        t.Errorf("SourcePort mismatch. Expected %d, got %d", udpHeader.SourcePort, deserialized.SourcePort)
    }

    if deserialized.DestinationPort != udpHeader.DestinationPort {
        t.Errorf("DestinationPort mismatch. Expected %d, got %d", udpHeader.DestinationPort, deserialized.DestinationPort)
    }

    if deserialized.Length != udpHeader.Length {
        t.Errorf("Length mismatch. Expected %d, got %d", udpHeader.Length, deserialized.Length)
    }

    if len(deserialized.Payload) != len(payload) {
        t.Errorf("Payload length mismatch. Expected %d, got %d", len(payload), len(deserialized.Payload))
    }

    if string(deserialized.Payload) != string(payload) {
        t.Errorf("Payload mismatch. Expected %s, got %s", payload, deserialized.Payload)
    }
}
