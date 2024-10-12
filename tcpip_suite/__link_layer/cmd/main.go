// cmd/main.go

package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "os"

    "link_layer/pkg/arp"
    "link_layer/pkg/ethernet"
    "link_layer/pkg/mtu"
    "link_layer/pkg/security"
    "link_layer/pkg/utils"
    "gopkg.in/yaml.v2"
)

type Config struct {
    LinkLayer struct {
        Interface string `yaml:"interface"`
        MTU       int    `yaml:"mtu"`
        Security  struct {
            MACsec struct {
                Enabled     bool   `yaml:"enabled"`
                CipherSuite string `yaml:"cipher_suite"`
                Key         string `yaml:"key"`
            } `yaml:"macsec"`
        } `yaml:"security"`
    } `yaml:"link_layer"`
}

func main() {
    // Parse command-line arguments for config file.
    configFile := flag.String("config", "configs/config.yaml", "Path to configuration file")
    flag.Parse()

    // Load configuration.
    configData, err := ioutil.ReadFile(*configFile)
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }

    var config Config
    err = yaml.Unmarshal(configData, &config)
    if err != nil {
        log.Fatalf("Failed to parse config file: %v", err)
    }

    // Initialize Logger.
    utils.InitLogger("link_layer.log")
    utils.Logger.Println("Starting Link Layer Implementation.")

    // Set MTU.
    err = mtu.SetMTU(config.LinkLayer.Interface, config.LinkLayer.MTU)
    if err != nil {
        utils.Logger.Fatalf("Failed to set MTU: %v", err)
    }

    // Initialize MACsec.
    macsecConfig := security.MACsec{
        Enabled:     config.LinkLayer.Security.MACsec.Enabled,
        CipherSuite: config.LinkLayer.Security.MACsec.CipherSuite,
        Key:         config.LinkLayer.Security.MACsec.Key,
    }
    err = macsecConfig.Initialize(config.LinkLayer.Interface)
    if err != nil {
        utils.Logger.Fatalf("Failed to initialize MACsec: %v", err)
    }

    // Example: Create and Serialize an Ethernet Frame.
    destMAC, err := net.ParseMAC("ff:ff:ff:ff:ff:ff") // Broadcast.
    if err != nil {
        utils.Logger.Fatalf("Invalid destination MAC: %v", err)
    }
    srcMAC, err := net.ParseMAC("00:0a:95:9d:68:16") // Example source MAC.
    if err != nil {
        utils.Logger.Fatalf("Invalid source MAC: %v", err)
    }
    payload := []byte("Hello, Ethernet!")
    ethFrame := ethernet.NewEthernetFrame(destMAC, srcMAC, ethernet.EtherTypeIPv4, payload)
    serializedFrame, err := ethFrame.Serialize()
    if err != nil {
        utils.Logger.Fatalf("Failed to serialize Ethernet frame: %v", err)
    }
    utils.Logger.Printf("Serialized Ethernet Frame: %x", serializedFrame)

    // Example: Create and Serialize an ARP Packet.
    arpRequest := arp.NewARPPacket(
        arp.ARPRequest,
        srcMAC,
        net.ParseIP("192.168.1.2"),
        net.ParseIP("192.168.1.3"),
        net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Target MAC unknown.
    )
    serializedARP, err := arpRequest.Serialize()
    if err != nil {
        utils.Logger.Fatalf("Failed to serialize ARP packet: %v", err)
    }
    utils.Logger.Printf("Serialized ARP Packet: %x", serializedARP)

    fmt.Println("Link Layer components initialized successfully.")
}