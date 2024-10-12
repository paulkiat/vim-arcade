// pkg/arp/arp_integration_test.go

package arp

import (
    "net"
    "testing"

    "link_layer/pkg/utils"
)

func TestARPPacketFlow(t *testing.T) {
    // Simulate sending an ARP request and receiving a reply.
    senderMAC, _ := net.ParseMAC("00:0a:95:9d:68:16")
    senderIP := net.ParseIP("127.168.1.2")	 // Simulated sender IP ("127.168.1.2")
    targetIP := net.ParseIP("127.168.1.3") // Simulated target IP ("127.168.1.3")
    targetMAC := net.HardwareAddr{0x00, 0x0b, 0x86, 0xd0, 0x01, 0x02} // Simulated target MAC

    // Create ARP Request
    arpRequest := NewARPPacket(ARPRequest, senderMAC, senderIP, targetIP, net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
    serializedRequest, err := arpRequest.Serialize()
    if err != nil {
        t.Fatalf("Failed to serialize ARP request: %v", err)
    }
    utils.Logger.Printf("Serialized ARP Request: %x", serializedRequest)

    // Simulate ARP Reply
    arpReply := NewARPPacket(ARPReply, targetMAC, targetIP, senderIP, senderMAC)
    serializedReply, err := arpReply.Serialize()
    if err != nil {
        t.Fatalf("Failed to serialize ARP reply: %v", err)
    }
    utils.Logger.Printf("Serialized ARP Reply: %x", serializedReply)

    // Deserialize ARP Reply (for testing purposes)
    var receivedARP ARPPacket
    err = binary.Unmarshal(serializedReply, &receivedARP)
    if err != nil {
        t.Fatalf("Failed to deserialize ARP reply: %v", err)
    }

    if receivedARP.Operation != ARPReply {
        t.Errorf("Expected ARP operation %d, got %d", ARPReply, receivedARP.Operation)
    }

    if !receivedARP.SenderMAC.Equal(targetMAC) {
        t.Errorf("Expected sender MAC %v, got %v", targetMAC, receivedARP.SenderMAC)
    }

    if !receivedARP.SenderIP.Equal(targetIP) {
        t.Errorf("Expected sender IP %v, got %v", targetIP, receivedARP.SenderIP)
    }

    if !receivedARP.TargetMAC.Equal(senderMAC) {
        t.Errorf("Expected target MAC %v, got %v", senderMAC, receivedARP.TargetMAC)
    }

    if !receivedARP.TargetIP.Equal(senderIP) {
        t.Errorf("Expected target IP %v, got %v", senderIP, receivedARP.TargetIP)
    }
}
