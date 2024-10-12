// Implements the Routing Information Protocol (RIP) for dynamic routing.

// dynamic_routing/main.go

package main

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "sync"
    "time"

    "../ip_routing"
    "../logger"
    "gopkg.in/yaml.v2"
)

type Config struct {
    DynamicRouting struct {
        Protocol  string   `yaml:"protocol"`
        Neighbors []string `yaml:"neighbors"`
    } `yaml:"dynamic_routing"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
    Server struct {
        IP string `yaml:"ip"`
    } `yaml:"server"`
}

type RIPPacket struct {
    Command  byte
    Version  byte
    Entries  []RIPEntry
}

type RIPEntry struct {
    Address  net.IP
    Subnet   net.IP
    NextHop  net.IP
    Metric   uint32
}

type RIPService struct {
    config      Config
    routing     *ip_routing.IPRouting
    neighbors   []net.IP
    conn        *net.UDPConn
    mutex       sync.Mutex
    logger      *log.Logger
}

func NewRIPService(config Config, routing *ip_routing.IPRouting, logger *log.Logger) *RIPService {
    neighbors := []net.IP{}
    for _, n := range config.DynamicRouting.Neighbors {
        ip := net.ParseIP(n)
        if ip != nil {
            neighbors = append(neighbors, ip)
        }
    }

    addr, err := net.ResolveUDPAddr("udp", ":520")
    if err != nil {
        logger.Fatalf("Failed to resolve UDP address for RIP: %v", err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        logger.Fatalf("Failed to listen on UDP port 520: %v", err)
    }

    return &RIPService{
        config:    config,
        routing:   routing,
        neighbors: neighbors,
        conn:      conn,
        logger:    logger,
    }
}

func (rs *RIPService) Start() {
    go rs.listen()
    go rs.advertise()
}

func (rs *RIPService) listen() {
    buffer := make([]byte, 512)
    for {
        n, addr, err := rs.conn.ReadFromUDP(buffer)
        if err != nil {
            rs.logger.Printf("Error reading UDP packet: %v", err)
            continue
        }

        packet, err := rs.parsePacket(buffer[:n])
        if err != nil {
            rs.logger.Printf("Failed to parse RIP packet from %s: %v", addr.IP.String(), err)
            continue
        }

        if packet.Version != 2 {
            rs.logger.Printf("Unsupported RIP version %d from %s", packet.Version, addr.IP.String())
            continue
        }

        if packet.Command == 1 { // Request
            rs.logger.Printf("Received RIP Request from %s", addr.IP.String())
            rs.sendResponse(addr.IP)
        } else if packet.Command == 2 { // Response
            rs.logger.Printf("Received RIP Response from %s", addr.IP.String())
            rs.updateRoutingTable(packet.Entries)
        }
    }
}

func (rs *RIPService) advertise() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        <-ticker.C
        rs.sendResponse("<broadcast>")
    }
}

func (rs *RIPService) sendResponse(target string) {
    rs.mutex.Lock()
    defer rs.mutex.Unlock()

    entries := []RIPEntry{}
    rs.routing.Mutex.Lock()
    for _, route := range rs.routing.RoutingTable {
        entry := RIPEntry{
            Address: net.ParseIP(route.Destination.IP.String()).To4(),
            Subnet:  net.ParseIP(route.Destination.Mask.String()).To4(),
            NextHop: route.Gateway.To4(),
            Metric:  1, // Directly connected
        }
        entries = append(entries, entry)
    }
    rs.routing.Mutex.Unlock()

    packet := RIPPacket{
        Command:  2, // Response
        Version:  2,
        Entries:  entries,
    }

    data, err := rs.buildPacket(packet)
    if err != nil {
        rs.logger.Printf("Failed to build RIP response packet: %v", err)
        return
    }

    var addr *net.UDPAddr
    if target == "<broadcast>" {
        addr, err = net.ResolveUDPAddr("udp", "224.0.0.9:520") // RIP multicast address
    } else {
        addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:520", target))
    }
    if err != nil {
        rs.logger.Printf("Failed to resolve RIP target address %s: %v", target, err)
        return
    }

    _, err = rs.conn.WriteToUDP(data, addr)
    if err != nil {
        rs.logger.Printf("Failed to send RIP response to %s: %v", addr.IP.String(), err)
        return
    }

    rs.logger.Printf("Sent RIP Response to %s", addr.IP.String())
}

func (rs *RIPService) parsePacket(data []byte) (*RIPPacket, error) {
    if len(data) < 4 {
        return nil, fmt.Errorf("packet too short")
    }

    packet := &RIPPacket{
        Command: data[0],
        Version: data[1],
    }

    numEntries := binary.BigEndian.Uint16(data[2:4])
    if len(data) < int(4+numEntries*16) {
        return nil, fmt.Errorf("packet too short for entries")
    }

    for i := 0; i < int(numEntries); i++ {
        start := 4 + i*16
        end := start + 16
        entryData := data[start:end]
        address := make(net.IP, 4)
        copy(address, entryData[0:4])
        subnet := make(net.IP, 4)
        copy(subnet, entryData[4:8])
        nextHop := make(net.IP, 4)
        copy(nextHop, entryData[8:12])
        metric := binary.BigEndian.Uint32(entryData[12:16])

        entry := RIPEntry{
            Address: address,
            Subnet:  subnet,
            NextHop: nextHop,
            Metric:  metric,
        }
        packet.Entries = append(packet.Entries, entry)
    }

    return packet, nil
}

func (rs *RIPService) buildPacket(packet RIPPacket) ([]byte, error) {
    buf := bytes.Buffer{}
    buf.WriteByte(packet.Command)
    buf.WriteByte(packet.Version)
    binary.Write(&buf, binary.BigEndian, uint16(len(packet.Entries)))
    for _, entry := range packet.Entries {
        buf.Write(entry.Address.To4())
        buf.Write(entry.Subnet.To4())
        buf.Write(entry.NextHop.To4())
        binary.Write(&buf, binary.BigEndian, entry.Metric)
    }
    return buf.Bytes(), nil
}

func (rs *RIPService) updateRoutingTable(entries []RIPEntry) {
    rs.routing.Mutex.Lock()
    defer rs.routing.Mutex.Unlock()

    for _, entry := range entries {
        destination := fmt.Sprintf("%s/%d", entry.Address.String(), maskToPrefix(entry.Subnet))
        existingRoute := rs.routing.FindRoute(destination)
        if existingRoute == nil || entry.Metric < existingRoute.Metric {
            rs.routing.AddRoute(destination, entry.NextHop.String(), "eth0")
            rs.logger.Printf("Updated route: %s via %s metric %d", destination, entry.NextHop.String(), entry.Metric)
        }
    }
}

func maskToPrefix(mask net.IP) int {
    ones, _ := net.IPMask(mask).Size()
    return ones
}

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
    loggerInstance := logger.SetupLogging(logger.Config{
        Logging: struct {
            Level string `yaml:"level"`
            File  string `yaml:"file"`
        }{
            Level: config.Logging.Level,
            File:  config.Logging.File,
        },
    })

    // Initialize IP Routing
    routing := ip_routing.NewIPRouting(config, loggerInstance)
    defer loggerInstance.Println("IP Routing service terminated.")

    // Initialize RIP Service
    ripService := NewRIPService(config, routing, loggerInstance)
    ripService.Start()
    defer ripService.conn.Close()

    // Command Interface
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("Dynamic Routing Service (RIP)")
    fmt.Println("Commands: add <destination> <gateway>, remove <destination>, display, quit")
    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        input := scanner.Text()
        parts := strings.Fields(input)
        if len(parts) == 0 {
            continue
        }

        switch parts[0] {
        case "add":
            if len(parts) != 3 {
                fmt.Println("Usage: add <destination> <gateway>")
                continue
            }
            destination := parts[1]
            gateway := parts[2]
            routing.AddRoute(destination, gateway, "eth0")
        case "remove":
            if len(parts) != 2 {
                fmt.Println("Usage: remove <destination>")
                continue
            }
            destination := parts[1]
            routing.RemoveRoute(destination)
        case "display":
            routing.DisplayRoutingTable()
        case "quit":
            return
        default:
            fmt.Println("Unknown command.")
        }
    }
}