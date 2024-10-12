// pkg/mtu/mtu_test.go

package mtu

import (
    "net"
    "testing"
)

func TestGetMTU(t *testing.T) {
    ifaceName := "lo" // Loopback interface
    mtu, err := GetMTU(ifaceName)
    if err != nil {
        t.Fatalf("Failed to get MTU for interface %s: %v", ifaceName, err)
    }

    if mtu <= 0 {
        t.Errorf("Invalid MTU value retrieved: %d", mtu)
    }
}

func TestSetMTU(t *testing.T) {
    ifaceName := "lo" // Loopback interface

    // Attempt to set MTU (requires elevated privileges)
    // This test may fail if not run as root.
    err := SetMTU(ifaceName, 1500)
    if err != nil {
        t.Fatalf("Failed to set MTU for interface %s: %v", ifaceName, err)
    }

    mtu, err := GetMTU(ifaceName)
    if err != nil {
        t.Fatalf("Failed to get MTU for interface %s: %v", ifaceName, err)
    }

    if mtu != 1500 {
        t.Errorf("MTU not set correctly. Expected 1500, got %d", mtu)
    }
}