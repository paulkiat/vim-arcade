// pkg/vlan/vlan_test.go

package vlan

import (
    "net"
    "testing"
)

func TestVLANAddRemove(t *testing.T) {
    vlanConfig := VLAN{
        Enabled:  true,
        VLANID:   100,
        Priority: 5,
    }

    ifaceName := "lo" // Loopback interface for testing.

    // Add VLAN
    err := vlanConfig.AddVLAN(ifaceName)
    if err != nil {
        t.Fatalf("Failed to add VLAN: %v", err)
    }

    // Remove VLAN
    err = vlanConfig.RemoveVLAN(ifaceName)
    if err != nil {
        t.Fatalf("Failed to remove VLAN: %v", err)
    }
}

func TestVLANFrameSerialization(t *testing.T) {
	destMAC, _ := net.ParseMAC("ff:ff:ff:ff:ff:ff")
	srcMAC, _ := net.ParseMAC("00:0a:95:9d:68:16")
	etherType := uint16(0x0800) // IPv4
	vlanID := uint16(10)
	payload := []byte("Test VLAN Payload")

	vlanFrame := NewVLANFrame(destMAC, srcMAC, etherType, vlanID, payload)
	serialized, err := vlanFrame.Serialize()
	if err != nil {
			t.Fatalf("Failed to serialize VLAN frame: %v", err)
	}

	deserialized, err := DeserializeVLANFrame(serialized)
	if err != nil {
			t.Fatalf("Failed to deserialize VLAN frame: %v", err)
	}

	if deserialized.VLANID != vlanID {
			t.Errorf("VLAN ID mismatch. Expected %d, got %d", vlanID, deserialized.VLANID)
	}

	if !deserialized.EthernetFrame.DestinationMAC.Equal(destMAC) {
			t.Errorf("Destination MAC mismatch. Expected %v, got %v", destMAC, deserialized.EthernetFrame.DestinationMAC)
	}

	if !deserialized.EthernetFrame.SourceMAC.Equal(srcMAC) {
			t.Errorf("Source MAC mismatch. Expected %v, got %v", srcMAC, deserialized.EthernetFrame.SourceMAC)
	}

	if deserialized.EthernetFrame.EtherType != etherType {
			t.Errorf("EtherType mismatch. Expected 0x%X, got 0x%X", etherType, deserialized.EthernetFrame.EtherType)
	}

	if string(deserialized.EthernetFrame.Payload) != string(payload) {
			t.Errorf("Payload mismatch. Expected %s, got %s", payload, deserialized.EthernetFrame.Payload)
	}
}