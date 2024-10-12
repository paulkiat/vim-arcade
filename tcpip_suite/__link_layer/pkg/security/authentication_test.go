// pkg/security/authentication_test.go

package security

import (
	"testing"
)

func TestAuthenticationInitialization(t *testing.T) {
	interfaceName := "lo" // Loopback interface (for testing purposes)

	authConfig := Authentication{
		Enabled:   false, // Disable Authentication for testing
		PortBased: true,
		EAPType:   "PEAP",
		Username:  "admin",
		Password:  "password",
	}

	err := authConfig.Initialize(interfaceName)
	if err != nil {
		t.Fatalf("Failed to initialize 802.1X Authentication: %v", err)
	}

	// Since Authentication is disabled, no further checks are needed.
}