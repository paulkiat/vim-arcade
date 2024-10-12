package logger

import (
    "log"
    "os"
    "strings"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
}

func SetupLogging(config Config) *log.Logger {
    var flag int
    switch strings.ToUpper(config.Logging.Level) {
    case "DEBUG":
        flag = log.Ldate | log.Ltime | log.Lshortfile
    case "INFO":
        flag = log.Ldate | log.Ltime
    case "WARNING":
        flag = log.Ldate | log.Ltime | log.Lshortfile
    case "ERROR":
        flag = log.Ldate | log.Ltime | log.Lshortfile
    case "CRITICAL":
        flag = log.Ldate | log.Ltime | log.Lshortfile
    default:
        flag = log.Ldate | log.Ltime
    }

    file, err := os.OpenFile(config.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file %s: %v", config.Logging.File, err)
    }

    logger := log.New(file, "", flag)
    return logger
}