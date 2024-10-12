// snmp_manager/logger.go

package snmp_manager

import (
    "io"
    "log"
    "os"
)

// Logger is the dedicated logger for the SNMP Manager module.
var Logger *log.Logger

// InitLogger initializes the logger for the SNMP Manager module.
func InitLogger(moduleName string, logFile string) {
    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file %s: %v", logFile, err)
    }

    multiWriter := io.MultiWriter(os.Stdout, file)
    Logger = log.New(multiWriter, moduleName+" - ", log.LstdFlags|log.Lshortfile)
}