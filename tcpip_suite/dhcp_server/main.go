// This module implements the DHCP server capable of handling DHCP Discover and Request messages and responding with DHCP Offer and Acknowledgment messages.

// dhcp_server/main.go

package main

import (
    "encoding/binary"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "time"

    "gopkg.in/yaml.v2"

    "../lease_manager"
    "../logger"
)

type Config struct {
    Server struct {
        Interface    string   `yaml:"interface"`
        IP           string   `yaml:"ip"`
        SubnetMask   string   `yaml:"subnet_mask"`
        Router       string   `yaml:"router"`
        DNSServers   []string `yaml:"dns_servers"`
        IPPool       struct {
            Start string `yaml:"start"`
            End   string `yaml:"end"`
        } `yaml:"ip_pool"`
        LeaseTime int `yaml:"lease_time"`
    } `yaml:"server"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
}

var loggerInstance *log.Logger
var leaseMgr *lease_manager.LeaseManager

func main() {
    // Load configuration
    configFile, err := os.Open("config.yaml")
    if err != nil {
        log.Fatalf("Failed to open config.yaml: %v", err)
    }
    defer configFile.Close()

    decoder := yaml.NewDecoder(configFile)
    var config Config
    err = decoder.Decode(&config)
    if err != nil {
        log.Fatalf("Failed to parse config.yaml: %v", err)
    }

    // Setup logging
    loggerInstance = logger.SetupLogging(logger.Config{
        Logging: struct {
            Level string `yaml:"level"`
            File  string `yaml:"file"`
        }{
            Level: config.Logging.Level,
            File:  config.Logging.File,
        },
    })

    // Initialize Lease Manager
    leaseMgr = lease_manager.NewLeaseManager(config.Server.IPPool.Start, config.Server.IPPool.End, config.Server.LeaseTime, loggerInstance)
    defer leaseMgr.Stop()

    // Setup UDP listener
    addr, err := net.ResolveUDPAddr("udp", ":67")
    if err != nil {
        loggerInstance.Fatalf("Failed to resolve UDP address: %v", err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        loggerInstance.Fatalf("Failed to listen on UDP port 67: %v", err)
    }
    defer conn.Close()

    loggerInstance.Println("DHCP Server is running...")

    buffer := make([]byte, 1024)
    for {
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            loggerInstance.Printf("Error reading UDP packet: %v", err)
            continue
        }

        go handlePacket(conn, buffer[:n], clientAddr, config)
    }
}

func handlePacket(conn *net.UDPConn, data []byte, clientAddr *net.UDPAddr, config Config) {
    messageType := getDHCPMessageType(data)
    clientMAC := getClientMAC(data)
    loggerInstance.Printf("Received DHCP Message Type: %d from MAC: %s", messageType, clientMAC)

    switch messageType {
    case 1: // DHCP Discover
        handleDiscover(conn, data, clientAddr, config)
    case 3: // DHCP Request
        handleRequest(conn, data, clientAddr, config)
    default:
        loggerInstance.Printf("Unhandled DHCP Message Type: %d from MAC: %s", messageType, clientMAC)
    }
}

func getDHCPMessageType(data []byte) byte {
    // DHCP options start at byte 240
    options := data[240:]
    for i := 0; i < len(options); i++ {
        if options[i] == 53 && i+2 < len(options) {
            return options[i+2]
        }
    }
    return 0
}

func getClientMAC(data []byte) string {
    mac := data[28:34]
    return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
        mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func handleDiscover(conn *net.UDPConn, data []byte, clientAddr *net.UDPAddr, config Config) {
    clientMAC := getClientMAC(data)
    offeredIP := leaseMgr.AssignIP(clientMAC)
    if offeredIP == "" {
        loggerInstance.Println("No available IP addresses to offer.")
        return
    }

    offerPacket := buildDHCPOffer(data, offeredIP, config)
    _, err := conn.WriteToUDP(offerPacket, &net.UDPAddr{
        IP:   net.IPv4bcast,
        Port: 68,
    })
    if err != nil {
        loggerInstance.Printf("Failed to send DHCP Offer to %s: %v", clientMAC, err)
        return
    }

    loggerInstance.Printf("Sent DHCP Offer to %s with IP %s", clientMAC, offeredIP)
}

func handleRequest(conn *net.UDPConn, data []byte, clientAddr *net.UDPAddr, config Config) {
    clientMAC := getClientMAC(data)
    requestedIP := getRequestedIP(data)
    loggerInstance.Printf("DHCP Request for IP %s from MAC %s", requestedIP, clientMAC)

    assignedIP := leaseMgr.GetIP(clientMAC)
    if assignedIP != requestedIP {
        loggerInstance.Printf("IP mismatch for MAC %s: requested %s, assigned %s", clientMAC, requestedIP, assignedIP)
        sendDHCPNAK(conn, data, config)
        return
    }

    ackPacket := buildDHCPAck(data, assignedIP, config)
    _, err := conn.WriteToUDP(ackPacket, &net.UDPAddr{
        IP:   net.IPv4bcast,
        Port: 68,
    })
    if err != nil {
        loggerInstance.Printf("Failed to send DHCP ACK to %s: %v", clientMAC, err)
        return
    }

    loggerInstance.Printf("Sent DHCP ACK to %s with IP %s", clientMAC, assignedIP)
}

func getRequestedIP(data []byte) string {
    options := data[240:]
    for i := 0; i < len(options); i++ {
        if options[i] == 50 && i+5 <= len(options) { // DHCP_OPTION_REQUESTED_IP has length 4
            ipBytes := options[i+2 : i+6]
            return net.IP(ipBytes).String()
        }
    }
    // Fallback to 'yiaddr' field if 'requested_ip' not found
    return net.IP(data[16:20]).String()
}

func sendDHCPNAK(conn *net.UDPConn, data []byte, config Config) {
    nakPacket := buildDHCPNAK(data, config)
    _, err := conn.WriteToUDP(nakPacket, &net.UDPAddr{
        IP:   net.IPv4bcast,
        Port: 68,
    })
    if err != nil {
        loggerInstance.Printf("Failed to send DHCP NAK: %v", err)
        return
    }

    clientMAC := getClientMAC(data)
    loggerInstance.Printf("Sent DHCP NAK to %s", clientMAC)
}

func buildDHCPOffer(data []byte, offeredIP string, config Config) []byte {
    offer := make([]byte, 300)

    // BOOTP fields
    offer[0] = 2  // op: BOOTREPLY
    offer[1] = 1  // htype: Ethernet
    offer[2] = 6  // hlen: MAC length
    offer[3] = 0  // hops
    xid := binary.BigEndian.Uint32(data[4:8])
    binary.BigEndian.PutUint32(offer[4:8], xid) // xid
    binary.BigEndian.PutUint16(offer[8:10], 0)  // secs
    binary.BigEndian.PutUint16(offer[10:12], 0x8000) // flags: Broadcast
    copy(offer[12:16], []byte{0, 0, 0, 0}) // ciaddr
    copy(offer[16:20], net.ParseIP(offeredIP).To4()) // yiaddr
    copy(offer[20:24], net.ParseIP(config.Server.IP).To4()) // siaddr
    copy(offer[24:28], []byte{0, 0, 0, 0}) // giaddr
    copy(offer[28:34], data[28:34]) // chaddr (MAC)
    // Rest of chaddr is padding

    // Magic cookie
    copy(offer[236:240], []byte{0x63, 0x82, 0x53, 0x63})

    // DHCP Options
    idx := 240
    offer[idx] = 53 // DHCP Message Type
    offer[idx+1] = 1  // Length
    offer[idx+2] = 2  // DHCPOFFER
    idx += 3

    offer[idx] = 54 // DHCP Server Identifier
    offer[idx+1] = 4
    copy(offer[idx+2:idx+6], net.ParseIP(config.Server.IP).To4())
    idx += 6

    offer[idx] = 1 // Subnet Mask
    offer[idx+1] = 4
    copy(offer[idx+2:idx+6], net.ParseIP(config.Server.SubnetMask).To4())
    idx += 6

    offer[idx] = 3 // Router
    offer[idx+1] = 4
    copy(offer[idx+2:idx+6], net.ParseIP(config.Server.Router).To4())
    idx += 6

    offer[idx] = 6 // DNS Servers
    dnsLen := len(config.Server.DNSServers) * 4
    offer[idx+1] = byte(dnsLen)
    for _, dns := range config.Server.DNSServers {
        copy(offer[idx+2:idx+6], net.ParseIP(dns).To4())
        idx += 4
    }

    offer[idx] = 51 // IP Lease Time
    offer[idx+1] = 4
    binary.BigEndian.PutUint32(offer[idx+2:idx+6], uint32(config.Server.LeaseTime))
    idx += 6

    offer[idx] = 255 // End Option
    // Pad remaining bytes with zeros
    return offer
}

func buildDHCPAck(data []byte, assignedIP string, config Config) []byte {
    ack := make([]byte, 300)

    // BOOTP fields
    ack[0] = 2  // op: BOOTREPLY
    ack[1] = 1  // htype: Ethernet
    ack[2] = 6  // hlen: MAC length
    ack[3] = 0  // hops
    xid := binary.BigEndian.Uint32(data[4:8])
    binary.BigEndian.PutUint32(ack[4:8], xid) // xid
    binary.BigEndian.PutUint16(ack[8:10], 0)  // secs
    binary.BigEndian.PutUint16(ack[10:12], 0x8000) // flags: Broadcast
    copy(ack[12:16], []byte{0, 0, 0, 0}) // ciaddr
    copy(ack[16:20], net.ParseIP(assignedIP).To4()) // yiaddr
    copy(ack[20:24], net.ParseIP(config.Server.IP).To4()) // siaddr
    copy(ack[24:28], []byte{0, 0, 0, 0}) // giaddr
    copy(ack[28:34], data[28:34]) // chaddr (MAC)
    // Rest of chaddr is padding

    // Magic cookie
    copy(ack[236:240], []byte{0x63, 0x82, 0x53, 0x63})

    // DHCP Options
    idx := 240
    ack[idx] = 53 // DHCP Message Type
    ack[idx+1] = 1  // Length
    ack[idx+2] = 5  // DHCPACK
    idx += 3

    ack[idx] = 54 // DHCP Server Identifier
    ack[idx+1] = 4
    copy(ack[idx+2:idx+6], net.ParseIP(config.Server.IP).To4())
    idx += 6

    ack[idx] = 1 // Subnet Mask
    ack[idx+1] = 4
    copy(ack[idx+2:idx+6], net.ParseIP(config.Server.SubnetMask).To4())
    idx += 6

    ack[idx] = 3 // Router
    ack[idx+1] = 4
    copy(ack[idx+2:idx+6], net.ParseIP(config.Server.Router).To4())
    idx += 6

    ack[idx] = 6 // DNS Servers
    dnsLen := len(config.Server.DNSServers) * 4
    ack[idx+1] = byte(dnsLen)
    for _, dns := range config.Server.DNSServers {
        copy(ack[idx+2:idx+6], net.ParseIP(dns).To4())
        idx += 4
    }

    ack[idx] = 51 // IP Lease Time
    ack[idx+1] = 4
    binary.BigEndian.PutUint32(ack[idx+2:idx+6], uint32(config.Server.LeaseTime))
    idx += 6

    ack[idx] = 255 // End Option
    // Pad remaining bytes with zeros
    return ack
}

func buildDHCPNAK(data []byte, config Config) []byte {
    nak := make([]byte, 300)

    // BOOTP fields
    nak[0] = 2  // op: BOOTREPLY
    nak[1] = 1  // htype: Ethernet
    nak[2] = 6  // hlen: MAC length
    nak[3] = 0  // hops
    xid := binary.BigEndian.Uint32(data[4:8])
    binary.BigEndian.PutUint32(nak[4:8], xid) // xid
    binary.BigEndian.PutUint16(nak[8:10], 0)  // secs
    binary.BigEndian.PutUint16(nak[10:12], 0x8000) // flags: Broadcast
    copy(nak[12:16], []byte{0, 0, 0, 0}) // ciaddr
    copy(nak[16:20], []byte{0, 0, 0, 0}) // yiaddr
    copy(nak[20:24], net.ParseIP(config.Server.IP).To4()) // siaddr
    copy(nak[24:28], []byte{0, 0, 0, 0}) // giaddr
    copy(nak[28:34], data[28:34]) // chaddr (MAC)
    // Rest of chaddr is padding

    // Magic cookie
    copy(nak[236:240], []byte{0x63, 0x82, 0x53, 0x63})

    // DHCP Options
    idx := 240
    nak[idx] = 53 // DHCP Message Type
    nak[idx+1] = 1  // Length
    nak[idx+2] = 6  // DHCPNAK
    idx += 3

    nak[idx] = 255 // End Option

    // Pad remaining bytes with zeros
    return nak
}