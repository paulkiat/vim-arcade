// Manages multicast group subscriptions and handles sending/receiving multicast messages.

// multicast/multicast_service.go

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
    Multicast struct {
        Groups []string `yaml:"groups"`
    } `yaml:"multicast"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
}

type MulticastService struct {
    groups  []net.IP
    conn    *net.UDPConn
    logger  *log.Logger
    mutex   sync.Mutex
}

func NewMulticastService(config Config, logger *log.Logger) *MulticastService {
    groups := []net.IP{}
    for _, g := range config.Multicast.Groups {
        ip := net.ParseIP(g)
        if ip != nil {
            groups = append(groups, ip)
        }
    }

    addr, err := net.ResolveUDPAddr("udp", ":5007") // Arbitrary port
    if err != nil {
        logger.Fatalf("Failed to resolve UDP address for multicast: %v", err)
    }

    conn, err := net.ListenMulticastUDP("udp", nil, addr)
    if err != nil {
        logger.Fatalf("Failed to listen on multicast UDP: %v", err)
    }

    return &MulticastService{
        groups: groups,
        conn:   conn,
        logger: logger,
    }
}

func (ms *MulticastService) Start() {
    go ms.listen()
}

func (ms *MulticastService) listen() {
    buffer := make([]byte, 1024)
    for {
        n, src, err := ms.conn.ReadFromUDP(buffer)
        if err != nil {
            ms.logger.Printf("Error reading multicast message: %v", err)
            continue
        }

        message := strings.TrimSpace(string(buffer[:n]))
        ms.logger.Printf("Received multicast message from %s: %s", src.IP.String(), message)
    }
}

func (ms *MulticastService) SendMessage(message string) {
    ms.mutex.Lock()
    defer ms.mutex.Unlock()

    for _, group := range ms.groups {
        addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:5007", group.String()))
        if err != nil {
            ms.logger.Printf("Failed to resolve multicast address %s: %v", group.String(), err)
            continue
        }

        _, err = ms.conn.WriteToUDP([]byte(message), addr)
        if err != nil {
            ms.logger.Printf("Failed to send multicast message to %s: %v", group.String(), err)
            continue
        }

        ms.logger.Printf("Sent multicast message to %s: %s", group.String(), message)
    }
}

func (ms *MulticastService) JoinGroup(group string) {
    ms.mutex.Lock()
    defer ms.mutex.Unlock()

    ip := net.ParseIP(group)
    if ip == nil {
        ms.logger.Printf("Invalid multicast group IP: %s", group)
        return
    }

    ms.groups = append(ms.groups, ip)
    ms.logger.Printf("Joined multicast group: %s", group)
}

func (ms *MulticastService) LeaveGroup(group string) {
    ms.mutex.Lock()
    defer ms.mutex.Unlock()

    index := -1
    for i, g := range ms.groups {
        if g.String() == group {
            index = i
            break
        }
    }

    if index != -1 {
        ms.groups = append(ms.groups[:index], ms.groups[index+1:]...)
        ms.logger.Printf("Left multicast group: %s", group)
    } else {
        ms.logger.Printf("Not a member of multicast group: %s", group)
    }
}

func (ms *MulticastService) Stop() {
    ms.conn.Close()
    ms.logger.Println("Multicast Service stopped.")
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

    // Initialize Multicast Service
    multicastService := NewMulticastService(config, loggerInstance)
    multicastService.Start()
    defer multicastService.Stop()

    // Command Interface
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("Multicast Service")
    fmt.Println("Commands: send <message>, join <group>, leave <group>, quit")
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
				case "send":
						if len(parts) < 2 {
								fmt.Println("Usage: send <message>")
								continue
						}
						multicastService.SendMessage(strings.Join(parts[1:], " "))
				case "join":
						if len(parts) < 2 {
								fmt.Println("Usage: join <group>")
								continue
						}
						multicastService.JoinGroup(parts[1])
				case "leave":
						if len(parts) < 2 {
								fmt.Println("Usage: leave <group>")
								continue
						}
						multicastService.LeaveGroup(parts[1])
				case "quit":
						return
				}
		}
}

func NewMulticastService(config Config, logger *log.Logger) *MulticastService {
		return &MulticastService{
				conn:    net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 5007}),
				logger:  logger,
				groups:  []net.IP{},
		}
}