// pkg/security/macsec.go

package security

import (
    "errors"
    "fmt"
    "os/exec"

    "link_layer/pkg/utils"
)

// MACsec represents a MACsec security configuration.
type MACsec struct {
    Enabled     bool
    CipherSuite string
    Key         string
}

// InitializeMACsec initializes MACsec on a given interface.
func (m *MACsec) Initialize(interfaceName string) error {
    if !m.Enabled {
        utils.Logger.Println("MACsec is disabled.")
        return nil
    }

    // Validate cipher suite.
    validSuites := map[string]bool{
        "GCM-AES-128": true,
        "GCM-AES-256": true,
    }

    if !validSuites[m.CipherSuite] {
        return errors.New("unsupported cipher suite for MACsec")
    }

    // Apply MACsec settings using system commands (Linux specific).
    // Requires iproute2 with MACsec support.
    // Example command:
    // ip link add link eth0 macsec0 type macsec cipher-suite GCM-AES-128 key "your-macsec-key-here"
    cmd := exec.Command("ip", "link", "add", "link", interfaceName, "macsec0", "type", "macsec", "cipher-suite", m.CipherSuite, "key", m.Key)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to add MACsec interface: %v, output: %s", err, string(output))
    }

    // Bring up the MACsec interface.
    cmd = exec.Command("ip", "link", "set", "macsec0", "up")
    output, err = cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to bring up MACsec interface: %v, output: %s", err, string(output))
    }

    utils.Logger.Printf("Initialized MACsec on interface %s with cipher suite %s", interfaceName, m.CipherSuite)
    return nil
}

// RemoveMACsec removes MACsec configuration from the interface.
func (m *MACsec) RemoveMACsec(interfaceName string) error {
    if !m.Enabled {
        utils.Logger.Println("MACsec is disabled.")
        return nil
    }

    // Delete the MACsec interface.
    cmd := exec.Command("ip", "link", "delete", "macsec0")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to delete MACsec interface: %v, output: %s", err, string(output))
    }

    utils.Logger.Printf("Removed MACsec from interface %s", interfaceName)
    return nil
}