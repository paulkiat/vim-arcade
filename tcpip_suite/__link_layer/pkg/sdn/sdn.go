// pkg/sdn/sdn.go

package sdn

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"

    "link_layer/pkg/utils"
)

// SDNController represents the SDN controller's API endpoints.
type SDNController struct {
    IP   				  string
    Port 				  int
		Enabled 		  bool
		ControllerURL string
}

// SDNConfig represents the configuration payload sent to the SDN controller.
type SDNConfig struct {
    Interface string `json:"interface"`
    MTU       int    `json:"mtu"`
    VLAN      struct {
        Enabled  bool `json:"enabled"`
        VLANID   int  `json:"vlan_id"`
        Priority int  `json:"priority"`
    } `json:"vlan"`
    Security struct {
        MACsec struct {
            Enabled     bool   `json:"enabled"`
            CipherSuite string `json:"cipher_suite"`
            Key         string `json:"key"`
        } `json:"macsec"`
        Auth8021X struct {
            Enabled      bool   `json:"enabled"`
            EAPMethod    string `json:"eap_method"`
            Identity     string `json:"identity"`
            Certificate  string `json:"certificate"`
            PrivateKey   string `json:"private_key"`
            CACertificate string `json:"ca_certificate"`
        } `json:"auth8021x"`
    } `json:"security"`
}

// ApplyConfig sends the configuration to the SDN controller.
func (sdn *SDNController) ApplyConfig(config SDNConfig) error {
    url := fmt.Sprintf("http://%s:%d/config", sdn.IP, sdn.Port)

		if !s.Enabled {
			utils.Logger.Println("SDN Integration is disabled.")
			return nil
		}

		config := FrameConfig{
			Interface: interfaceName,
			MTU:       mtu,
			VLANID:    vlanID,
		}

    jsonData, err := json.Marshal(config)
    if err != nil {
			return fmt.Errorf("failed to marshal SDN configuration: %v", err)
    }

    resp, err := http.Post(fmt.Sprintf("%s/link_layer/configure", s.ControllerURL), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to send configuration to SDN controller: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("SDN controller returned non-OK status: %s", resp.Status)
    }

		utils.Logger.Printf("Applied Link Layer configuration to SDN controller at %s", s.ControllerURL)
    return nil
}

// Example usage: Apply configuration to SDN controller.
func (sdn *SDNController) ExampleApply(config SDNConfig) error {
    return sdn.ApplyConfig(config)
}