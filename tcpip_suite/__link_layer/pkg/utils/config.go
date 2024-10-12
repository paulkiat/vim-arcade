// pkg/utils/config.go

package utils

import (
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

// Config represents the structure of the configuration file.
type Config struct {
    LinkLayer struct {
        Interface string `yaml:"interface" validate:"required"`
        MTU       int    `yaml:"mtu" validate:"min=68,max=9000"`
        VLAN      struct {
            Enabled  bool   `yaml:"enabled"`
            VLANID   int    `yaml:"vlan_id" validate:"min=1,max=4094"`
            Priority int    `yaml:"priority" validate:"min=0,max=7"`
        } `yaml:"vlan"`
        Security struct {
            MACsec struct {
                Enabled     bool   `yaml:"enabled"`
                CipherSuite string `yaml:"cipher_suite" validate:"oneof=GCM-AES-128 GCM-AES-256"`
                Key         string `yaml:"key" validate:"required,min=16,max=64"`
            } `yaml:"macsec"`
						Dot1X struct {
								Enabled  bool   `yaml:"enabled"`
								Username string `yaml:"username"`
								Password string `yaml:"password"`
					} `yaml:"dot1x"`
            Auth8021X struct {
                Enabled      bool   `yaml:"enabled"`
                EAPMethod    string `yaml:"eap_method" validate:"required,oneof=PEAP EAP-TLS"`
                Identity     string `yaml:"identity" validate:"required_if=EAPMethod PEAP"`
                Certificate  string `yaml:"certificate"` // Path to certificate for EAP-TLS
                PrivateKey   string `yaml:"private_key"` // Path to private key for EAP-TLS
                CA_Certificate string `yaml:"ca_certificate"` // Path to CA certificate
            } `yaml:"auth8021x"`
        } `yaml:"security"`
        SDN struct {
            Enabled      bool   `yaml:"enabled"`
            ControllerURL string `yaml:"controller_url" validate:"required_if=Enabled true,url"`
            ControllerIP string `yaml:"controller_ip" validate:"ip"`
            ControllerPort int  `yaml:"controller_port" validate:"min=1,max=65535"`
        } `yaml:"sdn"`
    } `yaml:"link_layer"`
}

// LoadConfig reads and parses the YAML configuration file.
func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
    }

    // Validation can be added here using a validation library
    validate := validator.New()
    err = validate.Struct(cfg)
    if err != nil {
        return nil, fmt.Errorf("configuration validation error: %v", err)
    }

    return &config, nil
}