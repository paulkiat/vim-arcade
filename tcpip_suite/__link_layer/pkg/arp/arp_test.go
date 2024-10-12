// pkg/arp/arp_test.go

package arp

import (
    "net"
    "testing"
)

func TestARPPacketSerialization(t *testing.T) {
    senderMAC, _ := net.ParseMAC("00:0a:95:9d:68:16")
    senderIP := net.ParseIP("192.168.1.2")
    targetIP := net.ParseIP("192.168.1.3")
    targetMAC := net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // Unknown

    arpPacket := NewARPPacket(ARPRequest, senderMAC, senderIP, targetIP, targetMAC)
    serialized, err := arpPacket.Serialize()
    if err != nil {
        t.Fatalf("Failed to serialize ARP packet: %v", err)
    }

    if len(serialized) != 28 {
        t.Errorf("Serialized ARP packet length mismatch. Expected 28, got %d", len(serialized))
    }

    // Deserialize and verify
    deserializedPacket, err := DeserializeARPPacket(serialized)
    if err != nil {
        t.Fatalf("Failed to deserialize ARP packet: %v", err)
    }

    if deserializedPacket.HardwareType != arpPacket.HardwareType {
        t.Errorf("HardwareType mismatch. Expected %v, got %v", arpPacket.HardwareType, deserializedPacket.HardwareType)
    }
    if deserializedPacket.ProtocolType != arpPacket.ProtocolType {
        t.Errorf("ProtocolType mismatch. Expected %v, got %v", arpPacket.ProtocolType, deserializedPacket.ProtocolType)
    }
    if deserializedPacket.Operation != arpPacket.Operation {
        t.Errorf("Operation mismatch. Expected %v, got %v", arpPacket.Operation, deserializedPacket.Operation)
    }
    if !deserializedPacket.SenderMAC.Equal(senderMAC) {
        t.Errorf("SenderMAC mismatch. Expected %v, got %v", senderMAC, deserializedPacket.SenderMAC)
    }
    if !deserializedPacket.SenderIP.Equal(senderIP) {
        t.Errorf("SenderIP mismatch. Expected %v, got %v", senderIP, deserializedPacket.SenderIP)
    }
    if !deserializedPacket.TargetMAC.Equal(targetMAC) {
        t.Errorf("TargetMAC mismatch. Expected %v, got %v", targetMAC, deserializedPacket.TargetMAC)
    }
    if !deserializedPacket.TargetIP.Equal(targetIP) {
        t.Errorf("TargetIP mismatch. Expected %v, got %v", targetIP, deserializedPacket.TargetIP)
    }
}