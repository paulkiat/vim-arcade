// traceroute/main.go

package main

import (
    "flag"
    "fmt"
    "os"

    "enterprise_networking_suite/logger"
    "enterprise_networking_suite/traceroute"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: traceroute [options] <destination>")
        fmt.Println("Options:")
        fmt.Println("  -m <maxHops>    Maximum number of hops (default: 30)")
        fmt.Println("  -w <timeout>    Timeout in seconds for each probe (default: 5)")
        fmt.Println("  -q <queries>    Number of probes per hop (default: 3)")
        fmt.Println("  -g <router>     Specify an intermediate router for source routing (loose)")
        fmt.Println("  -G <router>     Specify an intermediate router for source routing (strict)")
        os.Exit(1)
    }

    // Initialize Logger
    logger.InitLogger("traceroute", "traceroute.log")

    // Parse flags
    maxHops := flag.Int("m", 30, "Maximum number of hops")
    timeout := flag.Int("w", 5, "Timeout in seconds for each probe")
    queries := flag.Int("q", 3, "Number of probes per hop")
    looseSourceRoutes := flag.String("g", "", "Specify an intermediate router for loose source routing (comma-separated)")
    strictSourceRoutes := flag.String("G", "", "Specify an intermediate router for strict source routing (comma-separated)")
    flag.Parse()

    if flag.NArg() < 1 {
        fmt.Println("Destination not specified.")
        os.Exit(1)
    }

    destination := flag.Arg(0)
    options := traceroute.TracerouteOptions{
        MaxHops:           *maxHops,
        Timeout:           *timeout,
        QueriesPerHop:     *queries,
        LooseSourceRoutes: parseRouters(*looseSourceRoutes),
        StrictSourceRoutes: parseRouters(*strictSourceRoutes),
    }

    tracer, err := traceroute.NewTraceroute(destination, options)
    if err != nil {
        logger.Logger.Fatalf("Failed to initialize traceroute: %v", err)
    }

    tracer.Run()
    tracer.Statistics()
}

func parseRouters(routerList string) []string {
    if routerList == "" {
        return nil
    }
    routers := []string{}
    for _, r := range flag.Args() {
        if r != "" {
            routers = append(routers, r)
        }
    }
    // Alternatively, split by comma if multiple routers are provided as a single string
    return routers
}