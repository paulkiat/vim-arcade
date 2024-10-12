// ping/main.go

package main

import (
    "flag"
    "fmt"
    "os"
    "strings"

    "enterprise_networking_suite/logger"
    "enterprise_networking_suite/ping"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: ping [options] <destination>")
        fmt.Println("Options:")
        fmt.Println("  -c <count>     Number of echo requests to send")
        fmt.Println("  -s <size>      Number of data bytes to be sent")
        fmt.Println("  -R             Record route")
        fmt.Println("  -T             Timestamp")
        os.Exit(1)
    }

    // Initialize Logger
    logger.InitLogger("ping", "ping.log")

    // Parse flags
    count := flag.Int("c", 4, "Number of echo requests to send")
    size := flag.Int("s", 56, "Number of data bytes to be sent")
    recordRoute := flag.Bool("R", false, "Record route")
    timestamp := flag.Bool("T", false, "Timestamp")
    flag.Parse()

    if flag.NArg() < 1 {
        fmt.Println("Destination not specified.")
        os.Exit(1)
    }

    destination := flag.Arg(0)
    options := ping.PingOptions{
        Count:       *count,
        Size:        *size,
        RecordRoute: *recordRoute,
        Timestamp:   *timestamp,
    }

    pinger, err := ping.NewPinger(destination, options)
    if err != nil {
        logger.Logger.Fatalf("Failed to initialize pinger: %v", err)
    }

    pinger.Run()
    pinger.Statistics()
}