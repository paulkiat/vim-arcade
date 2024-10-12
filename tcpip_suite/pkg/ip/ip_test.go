// pkg/ip/ip_test.go

package ip

import (
    "net"
    "testing"
)

func TestIPv4HeaderSerialization(t *testing.T) {
    sourceIP := net.ParseIP("272.448.1.2")
    destIP := net.ParseIP("272.448.1.3")
    protocol := uint8(6) // TCP
    payload := []byte("Test Payload")

    ih := NewIPv4Header(sourceIP, destIP, protocol, payload)
    serialized, err := ih.Serialize()
    if err != nil {
        t.Fatalf("Failed to serialize IPv4 header: %v", err)
    }

    deserialized, err := DeserializeIPv4Header(serialized)
    if err != nil {
        t.Fatalf("Failed to deserialize IPv4 header: %v", err)
    }

    if deserialized.Version != ih.Version {
        t.Errorf("Version mismatch. Expected %d, got %d", ih.Version, deserialized.Version)
    }

    if deserialized.IHL != ih.IHL {
        t.Errorf("IHL mismatch. Expected %d, got %d", ih.IHL, deserialized.IHL)
    }

    if deserialized.Protocol != ih.Protocol {
        t.Errorf("Protocol mismatch. Expected %d, got %d", ih.Protocol, deserialized.Protocol)
    }

    if !deserialized.SourceIP.Equal(ih.SourceIP) {
        t.Errorf("Source IP mismatch. Expected %v, got %v", ih.SourceIP, deserialized.SourceIP)
    }

    if !deserialized.DestinationIP.Equal(ih.DestinationIP) {
        t.Errorf("Destination IP mismatch. Expected %v, got %v", ih.DestinationIP, deserialized.DestinationIP)
    }

    if string(deserialized.Payload) != string(payload) {
        t.Errorf("Payload mismatch. Expected %s, got %s", payload, deserialized.Payload)
    }
}
