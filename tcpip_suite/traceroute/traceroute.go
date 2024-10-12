// traceroute/traceroute.go

package traceroute

import (
    "fmt"
    "net"
    "os"
    "strings"
    "time"

    "enterprise_networking_suite/logger"

    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
)

type TracerouteOptions struct {
    MaxHops             int
    Timeout             int
    QueriesPerHop       int
    LooseSourceRoutes  []string
    StrictSourceRoutes []string
}

type Traceroute struct {
    Destination string
    Options     TracerouteOptions
    Conn        *icmp.PacketConn
    ID          int
    Port        int
    Stats       Statistics
}

type Statistics struct {
    TotalProbes     int
    ReceivedProbes  int
    LostProbes      int
    Hops            int
    HopDetails      []HopDetail
}

type HopDetail struct {
    Hop        int
    Addresses  []string
    RTTs       []time.Duration
    TimedOut   bool
}

func NewTraceroute(destination string, options TracerouteOptions) (*Traceroute, error) {
    conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
        return nil, fmt.Errorf("failed to listen for ICMP packets: %v", err)
    }

    id := os.Getpid() & 0xffff
    port := 33435 // Starting UDP port for traceroute

    return &Traceroute{
        Destination: destination,
        Options:     options,
        Conn:        conn,
        ID:          id,
        Port:        port,
        Stats:       Statistics{HopDetails: []HopDetail{}},
    }, nil
}

func (t *Traceroute) Run() {
    defer t.Conn.Close()

    destAddr, err := net.ResolveIPAddr("ip4", t.Destination)
    if err != nil {
        logger.Logger.Fatalf("Failed to resolve destination address: %v", err)
    }

    fmt.Printf("traceroute to %s (%s), %d hops max, 40 byte packets\n",
        t.Destination, destAddr.String(), t.Options.MaxHops)

    for ttl := 1; ttl <= t.Options.MaxHops; ttl++ {
        hopDetail := HopDetail{
            Hop:       ttl,
            Addresses: []string{},
            RTTs:      []time.Duration{},
            TimedOut:  true,
        }

        for probe := 0; probe < t.Options.QueriesPerHop; probe++ {
            // Set TTL
            p := ipv4.NewPacketConn(t.Conn)
            err := p.SetTTL(ttl)
            if err != nil {
                logger.Logger.Printf("Failed to set TTL: %v", err)
                continue
            }

            // Construct UDP Datagram
            udpAddr := &net.UDPAddr{
                IP:   destAddr.IP,
                Port: t.Port,
            }

            // Send UDP Datagram
            start := time.Now()
            _, err = t.Conn.WriteTo([]byte("Traceroute Probe"), udpAddr)
            if err != nil {
                logger.Logger.Printf("Failed to send UDP datagram: %v", err)
                continue
            }

            // Set Read Deadline
            t.Conn.SetReadDeadline(time.Now().Add(time.Duration(t.Options.Timeout) * time.Second))

            // Listen for ICMP Replies
            reply := make([]byte, 1500)
            n, peer, err := t.Conn.ReadFrom(reply)
            rtt := time.Since(start)

            if err != nil {
                // Timeout or other error
                hopDetail.RTTs = append(hopDetail.RTTs, time.Duration(0))
                continue
            }

            // Parse ICMP Message
            rm, err := icmp.ParseMessage(ipv4.ICMPTypeTimeExceeded.Protocol(), reply[:n])
            if err != nil {
                logger.Logger.Printf("Failed to parse ICMP message: %v", err)
                continue
            }

            switch rm.Type {
            case ipv4.ICMPTypeTimeExceeded:
                hopDetail.TimedOut = false
                hopDetail.Addresses = appendIfMissing(hopDetail.Addresses, peer.String())
                hopDetail.RTTs = append(hopDetail.RTTs, rtt)
            case ipv4.ICMPTypeDestinationUnreachable:
                if body, ok := rm.Body.(*icmp.DstUnreach); ok {
                    if body.Code == 3 { // Port Unreachable
                        hopDetail.TimedOut = false
                        hopDetail.Addresses = appendIfMissing(hopDetail.Addresses, peer.String())
                        hopDetail.RTTs = append(hopDetail.RTTs, rtt)
                        // Destination reached
                        t.Stats.Hops = ttl
                        t.Stats.TotalProbes += probe + 1
                        t.Stats.ReceivedProbes += probe + 1
                        t.Stats.HopDetails = append(t.Stats.HopDetails, hopDetail)
                        t.Port++ // Increment port for next traceroute run
                        return
                    }
                }
            default:
                // Unexpected ICMP type
                logger.Logger.Printf("Received unexpected ICMP message type: %v", rm.Type)
            }
        }

        t.Stats.TotalProbes += t.Options.QueriesPerHop
        t.Stats.HopDetails = append(t.Stats.HopDetails, hopDetail)
        fmt.Printf("%d ", ttl)
        if len(hopDetail.Addresses) == 0 {
            fmt.Printf("* ")
            for i := 1; i < t.Options.QueriesPerHop; i++ {
                fmt.Printf("* ")
            }
            fmt.Println()
        } else {
            for _, addr := range hopDetail.Addresses {
                names, err := net.LookupAddr(addr)
                var name string
                if err == nil && len(names) > 0 {
                    name = strings.TrimSuffix(names[0], ".")
                } else {
                    name = addr
                }
                fmt.Printf("%s (%s) ", name, addr)
            }
            for _, rtt := range hopDetail.RTTs {
                if rtt > 0 {
                    fmt.Printf("%.0f ms ", float64(rtt.Microseconds())/1000.0)
                } else {
                    fmt.Printf("* ")
                }
            }
            fmt.Println()
            t.Stats.ReceivedProbes += len(hopDetail.RTTs)
        }
    }
}

func appendIfMissing(slice []string, elem string) []string {
    for _, a := range slice {
        if a == elem {
            return slice
        }
    }
    return append(slice, elem)
}

func (t *Traceroute) Statistics() {
    fmt.Printf("\n--- %s traceroute statistics ---\n", t.Destination)
    fmt.Printf("%d probes transmitted, %d probes received, %.0f%% packet loss\n",
        t.Stats.TotalProbes, t.Stats.ReceivedProbes,
        float64(t.Stats.TotalProbes-t.Stats.ReceivedProbes)/float64(t.Stats.TotalProbes)*100)
    // Additional statistics can be added here
}