// pkg/security/macsec_test.go

package security

import (
    "testing"
)

func TestMACsecInitialization(t *testing.T) {
    macsecConfig := MACsec{
        Enabled:     false, // Disable MACsec for testing.
        CipherSuite: "GCM-AES-128",
        Key:         "0123456789abcdef", // Example key.
    }

    ifaceName := "lo" // Loopback interface (for testing purposes).

    err := macsecConfig.Initialize(ifaceName)
    if err != nil {
        t.Fatalf("Failed to initialize MACsec: %v", err)
    }

    // Since MACsec is disabled, no further checks are needed.
}

func TestMACsecEnableDisable(t *testing.T) {
    macsecConfig := MACsec{
        Enabled:     true,
        CipherSuite: "GCM-AES-128",
        Key:         "0123456789abcdef", // Example key.
    }

    ifaceName := "eth0" // Replace with a valid interface name for testing.

    err := macsecConfig.Initialize(ifaceName)
    if err != nil {
        t.Fatalf("Failed to initialize MACsec: %v", err)
    }

    err = macsecConfig.RemoveMACsec(ifaceName)
    if err != nil {
        t.Fatalf("Failed to remove MACsec: %v", err)
    }
}