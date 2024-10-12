// Manages static routing tables and forwards packets based on destination IP addresses.

// ip_routing/main.go

package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "sync"

    "../logger"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Routing struct {
        StaticRoutes []struct {
            Destination string `yaml:"destination"`
            Gateway     string `yaml:"gateway"`
        } `yaml:"static_routes"`
    } `yaml:"routing"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
}

type RoutingTableEntry struct {
    Destination net.IPNet
    Gateway     net.IP
    Interface   string
}

type IPRouting struct {
    routingTable []RoutingTableEntry
    mutex        sync.Mutex
    logger       *log.Logger
}

func NewIPRouting(config Config, logger *log.Logger) *IPRouting {
    ir := &IPRouting{
        routingTable: []RoutingTableEntry{},
        logger:       logger,
    }

    for _, route := range config.Routing.StaticRoutes {
        _, network, err := net.ParseCIDR(route.Destination)
        if err != nil {
            ir.logger.Printf("Invalid destination network: %s", route.Destination)
            continue
        }

        gateway := net.ParseIP(route.Gateway)
        if gateway == nil {
            ir.logger.Printf("Invalid gateway IP: %s", route.Gateway)
            continue
        }

        entry := RoutingTableEntry{
            Destination: *network,
            Gateway:     gateway,
            Interface:   "eth0", // Assuming eth0; can be made configurable
        }
        ir.routingTable = append(ir.routingTable, entry)
        ir.logger.Printf("Added static route: %s via %s on %s", network.String(), gateway.String(), entry.Interface)
    }

    return ir
}

func (ir *IPRouting) AddRoute(destination, gateway, iface string) {
    _, network, err := net.ParseCIDR(destination)
    if err != nil {
        ir.logger.Printf("Invalid destination network: %s", destination)
        return
    }

    gw := net.ParseIP(gateway)
    if gw == nil {
        ir.logger.Printf("Invalid gateway IP: %s", gateway)
        return
    }

    entry := RoutingTableEntry{
        Destination: *network,
        Gateway:     gw,
        Interface:   iface,
    }

    ir.mutex.Lock()
    defer ir.mutex.Unlock()
    ir.routingTable = append(ir.routingTable, entry)
    ir.logger.Printf("Added route: %s via %s on %s", network.String(), gw.String(), iface)
}

func (ir *IPRouting) RemoveRoute(destination string) {
    _, network, err := net.ParseCIDR(destination)
    if err != nil {
        ir.logger.Printf("Invalid destination network: %s", destination)
        return
    }

    ir.mutex.Lock()
    defer ir.mutex.Unlock()

    for i, route := range ir.routingTable {
        if route.Destination.String() == network.String() {
            ir.routingTable = append(ir.routingTable[:i], ir.routingTable[i+1:]...)
            ir.logger.Printf("Removed route: %s", network.String())
            return
        }
    }

    ir.logger.Printf("Route not found: %s", network.String())
}

func (ir *IPRouting) DisplayRoutingTable() {
    ir.mutex.Lock()
    defer ir.mutex.Unlock()

    fmt.Println("Current IP Routing Table:")
    for _, route := range ir.routingTable {
        fmt.Printf("%s via %s on %s\n", route.Destination.String(), route.Gateway.String(), route.Interface)
    }
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
    ipRouting := NewIPRouting(config, loggerInstance)
    defer loggerInstance.Println("IP Routing service terminated.")

    // Command Interface
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("IP Routing Service")
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
            ipRouting.AddRoute(parts[1], parts[2], "eth0")
        case "remove":
            if len(parts) != 2 {
                fmt.Println("Usage: remove <destination>")
                continue
            }
            ipRouting.RemoveRoute(parts[1])
        case "display":
            ipRouting.DisplayRoutingTable()
        case "quit":
            return
        default:
            fmt.Println("Unknown command.")
        }
    }
}
