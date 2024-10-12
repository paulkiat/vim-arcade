// pkg/ethernet/ethernet_test.go

package ethernet

import (
    "net"
    "testing"
)

func TestEthernetFrameSerialization(t *testing.T) {
    destMAC, _ := net.ParseMAC("ff:ff:ff:ff:ff:ff") // Broadcast
    srcMAC, _ := net.ParseMAC("00:0a:95:9d:68:16") // Example source MAC
    payload := []byte("Test Payload")
    frame := NewEthernetFrame(destMAC, srcMAC, EtherTypeIPv4, payload)

    serialized, err := frame.Serialize()
    if err != nil {
        t.Fatalf("Failed to serialize Ethernet frame: %v", err)
    }

    expectedLength := 14 + len(payload) + 4
    if len(serialized) != expectedLength {
        t.Errorf("Serialized frame length mismatch. Expected %d, got %d", expectedLength, len(serialized))
    }

		// Additional assertions can be added to verify frame content


    // Deserialize and verify
    deserializedFrame, err := DeserializeEthernetFrame(serialized)
    if err != nil {
        t.Fatalf("Failed to deserialize Ethernet frame: %v", err)
    }

    if !deserializedFrame.DestinationMAC.Equal(destMAC) {
        t.Errorf("Destination MAC mismatch. Expected %v, got %v", destMAC, deserializedFrame.DestinationMAC)
    }
    if !deserializedFrame.SourceMAC.Equal(srcMAC) {
        t.Errorf("Source MAC mismatch. Expected %v, got %v", srcMAC, deserializedFrame.SourceMAC)
    }
    if deserializedFrame.EtherType != frame.EtherType {
        t.Errorf("EtherType mismatch. Expected %v, got %v", frame.EtherType, deserializedFrame.EtherType)
    }
    if string(deserializedFrame.Payload) != string(payload) {
        t.Errorf("Payload mismatch. Expected %s, got %s", payload, deserializedFrame.Payload)
    }
}

func TestAddVLANTag(t *testing.T) {
	destMAC, _ := net.ParseMAC("ff:ff:ff:ff:ff:ff")
	srcMAC, _ := net.ParseMAC("00:0a:95:9d:68:16")
	payload := []byte("VLAN Payload")
	frame := NewEthernetFrame(destMAC, srcMAC, EtherTypeIPv4, payload)

	vlanID := uint16(100)
	err := frame.AddVLANTag(vlanID)
	if err != nil {
			t.Fatalf("Failed to add VLAN tag: %v", err)
	}

	if frame.EtherType != EtherType8021Q {
			t.Errorf("VLAN tag not added correctly. Expected EtherType 0x8100, got 0x%04x", frame.EtherType)
	}

	// Check the serialized frame
	serialized, err := frame.Serialize()
	if err != nil {
			t.Fatalf("Failed to serialize Ethernet frame with VLAN: %v", err)
	}

	// Verify VLAN tag
	expectedVLANType := uint16(0x8100)
	actualVLANType := binary.BigEndian.Uint16(serialized[12:14])
	if actualVLANType != expectedVLANType {
			t.Errorf("VLAN type mismatch. Expected 0x8100, got 0x%04x", actualVLANType)
	}

	actualVLANID := binary.BigEndian.Uint16(serialized[14:16])
	if actualVLANID != vlanID {
			t.Errorf("VLAN ID mismatch. Expected %d, got %d", vlanID, actualVLANID)
	}
}
