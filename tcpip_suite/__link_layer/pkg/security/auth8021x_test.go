// pkg/security/auth8021x_test.go

package security

import (
    "testing"
)

func TestAuth8021XInitialization(t *testing.T) {
    authConfig := Auth8021X{
        Enabled:       false, // Disable for testing.
        EAPMethod:     "PEAP",
        Identity:      "testuser",
        Certificate:   "/path/to/client.crt",
        PrivateKey:    "/path/to/client.key",
        CACertificate: "/path/to/ca.crt",
    }

    ifaceName := "lo" // Loopback interface (for testing purposes).

    err := authConfig.Initialize(ifaceName)
    if err != nil {
        t.Fatalf("Failed to initialize 802.1X Authentication: %v", err)
    }

    // Since 802.1X is disabled, no further checks are needed.
}

func TestAuth8021XEnableDisable(t *testing.T) {
    authConfig := Auth8021X{
        Enabled:       true,
        EAPMethod:     "EAP-TLS",
        Identity:      "testuser",
        Certificate:   "/path/to/client.crt",
        PrivateKey:    "/path/to/client.key",
        CACertificate: "/path/to/ca.crt",
    }

    ifaceName := "eth0" // Replace with a valid interface name for testing.

    err := authConfig.Initialize(ifaceName)
    if err != nil {
        t.Fatalf("Failed to initialize 802.1X Authentication: %v", err)
    }

    err = authConfig.Remove8021X(ifaceName)
    if err != nil {
        t.Fatalf("Failed to remove 802.1X Authentication: %v", err)
    }
}