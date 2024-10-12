// pkg/security/authentication.go

package security

import (
	"errors"
	"fmt"
	"os/exec"

	"link_layer/pkg/utils"
)

// Authentication represents an 802.1X authentication configuration.
type Authentication struct {
	Enabled   bool
	PortBased bool
	EAPType   string
	Username  string
	Password  string
}

// InitializeAuthentication initializes 802.1X authentication on a given interface.
func (a *Authentication) Initialize(interfaceName string) error {
	if !a.Enabled {
		utils.Logger.Println("802.1X Authentication is disabled.")
		return nil
	}

	// Validate EAP type.
	validEAPTypes := map[string]bool{
		"PEAP":     true,
		"EAP-TLS":  true,
		"TTLS":     true,
		"LEAP":     true,
		"FAST":     true,
		"EAP-MSCHAPV2": true,
	}

	if !validEAPTypes[a.EAPType] {
		return errors.New("unsupported EAP type for 802.1X Authentication")
	}

	// Configure authentication using system commands (Linux specific).
	// Example using wpa_supplicant for 802.1X.
	// Create a wpa_supplicant configuration file.
	configContent := fmt.Sprintf(`
network={
    key_mgmt=IEEE8021X
    eap=%s
    identity="%s"
    password="%s"
    phase1="peapver=0"
}
`, a.EAPType, a.Username, a.Password)

	configFilePath := "/etc/wpa_supplicant/wpa_supplicant.conf"

	// Write the configuration file.
	err := writeConfigFile(configFilePath, configContent)
	if err != nil {
		return fmt.Errorf("failed to write wpa_supplicant config: %v", err)
	}

	// Start wpa_supplicant with the configuration.
	cmd := exec.Command("wpa_supplicant", "-B", "-i", interfaceName, "-c", configFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start wpa_supplicant: %v, output: %s", err, string(output))
	}

	utils.Logger.Printf("Initialized 802.1X Authentication on interface %s with EAP type %s", interfaceName, a.EAPType)
	return nil
}

// writeConfigFile writes the given content to the specified file path.
func writeConfigFile(filePath, content string) error {
	// Implement file writing with appropriate permissions.
	// Placeholder implementation.
	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | sudo tee %s", content, filePath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to write config file: %v, output: %s", err, string(output))
	}
	return nil
}