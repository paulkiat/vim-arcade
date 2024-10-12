// snmp_manager/main.go

package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"

    "enterprise_networking_suite/snmp_manager"
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

type Config struct {
    SNMPManager snmp_manager.Config `yaml:"snmp_manager"`
}

func main() {
    // Parse command-line arguments for config file
    configFile := flag.String("config", "config.yaml", "Path to configuration file")
    flag.Parse()

    // Load configuration
    configData, err := ioutil.ReadFile(*configFile)
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }

    var config Config
    err = yaml.Unmarshal(configData, &config)
    if err != nil {
        log.Fatalf("Failed to parse config file: %v", err)
    }

    // Initialize Logger
    snmp_manager.InitLogger("SNMP_Manager", config.SNMPManager.Logging.File)

    // Initialize SNMP Manager
    manager, err := snmp_manager.NewManager(config.SNMPManager)
    if err != nil {
        snmp_manager.Logger.Fatalf("Failed to initialize SNMP Manager: %v", err)
    }

    // Start SNMP Manager operations
    go manager.Run()

    // Start Trap Handler
    go manager.HandleTraps()

    // Wait for termination signal (e.g., Ctrl+C)
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    <-sigChan

    // Gracefully stop the manager
    manager.Stop()
}