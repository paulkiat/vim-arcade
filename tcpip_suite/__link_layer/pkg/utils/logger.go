// pkg/utils/logger.go

package utils

import (
    "io"
    "log"
    "os"
    "github.com/sirupsen/logrus"
)

// Logger is the dedicated logger for the Link Layer module.
var Logger *logrus.Logger

// InitLogger initializes the logger for the Link Layer module.
func InitLogger(logFile string) {
    Logger = logrus.New()
    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file %s: %v", logFile, err)
    }

    multiWriter := io.MultiWriter(os.Stdout, file)
    Logger.SetOutput(multiWriter)
    Logger.SetFormatter(&logrus.JSONFormatter{})
    Logger.SetLevel(logrus.InfoLevel)
}