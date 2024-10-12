// ping/ping.go

package ping

import (
    "fmt"
    "net"
    "os"
    "time"

    "enterprise_networking_suite/logger"

    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
)

type PingOptions struct {
    Count       int
    Size        int
    RecordRoute bool
    Timestamp   bool
}

type Statistics struct {
    Transmitted int
    Received    int
    Lost        int
    MinRTT      time.Duration
    AvgRTT      time.Duration
    MaxRTT      time.Duration
}

type Pinger struct {
    Destination string
    Options     PingOptions
    Conn        *icmp.PacketConn
    ID          int
    Seq         int
    Stats       Statistics
    RTTs        []time.Duration
}

func NewPinger(destination string, options PingOptions) (*Pinger, error) {
    conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
        return nil, fmt.Errorf("failed to listen for ICMP packets: %v", err)
    }

    id := os.Getpid() & 0xffff

    return &Pinger{
        Destination: destination,
        Options:     options,
        Conn:        conn,
        ID:          id,
        Seq:         0,
        Stats: Statistics{
            Transmitted: 0,
            Received:    0,
            Lost:        0,
            MinRTT:      time.Duration(1<<63 - 1),
            AvgRTT:      0,
            MaxRTT:      0,
        },
        RTTs: make([]time.Duration, 0),
    }, nil
}

func (p *Pinger) Run() {
    defer p.Conn.Close()

    addr, err := net.ResolveIPAddr("ip4", p.Destination)
    if err != nil {
        logger.Logger.Fatalf("Failed to resolve destination address: %v", err)
    }

    for p.Stats.Transmitted < p.Options.Count {
        p.Stats.Transmitted++
        p.Seq++

        message := p.createICMPEchoRequest()

        start := time.Now()
        _, err := p.Conn.WriteTo(message, addr)
        if err != nil {
            logger.Logger.Printf("Failed to send ICMP echo request: %v", err)
            continue
        }

        p.Conn.SetReadDeadline(time.Now().Add(2 * time.Second))
        reply := make([]byte, 1500)
        n, peer, err := p.Conn.ReadFrom(reply)
        if err != nil {
            logger.Logger.Printf("Request timeout for icmp_seq %d", p.Seq)
            continue
        }
        duration := time.Since(start)
        p.RTTs = append(p.RTTs, duration)
        p.Stats.Received++

        rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
        if err != nil {
            logger.Logger.Printf("Failed to parse ICMP message: %v", err)
            continue
        }

        switch rm.Type {
        case ipv4.ICMPTypeEchoReply:
            echo, ok := rm.Body.(*icmp.Echo)
            if !ok {
                logger.Logger.Printf("Invalid ICMP echo reply body")
                continue
            }
            if echo.ID != p.ID || echo.Seq != p.Seq {
                logger.Logger.Printf("Received unmatched echo reply: ID=%d, Seq=%d", echo.ID, echo.Seq)
                continue
            }
            ttl := p.getTTL(peer)
            fmt.Printf("64 bytes from %s: icmp_seq=%d ttl=%d time=%.3f ms\n",
                addr.String(), p.Seq, ttl, float64(duration.Microseconds())/1000.0)
        default:
            logger.Logger.Printf("Received unexpected ICMP message: %v", rm.Type)
        }

        time.Sleep(1 * time.Second)
    }
}

func (p *Pinger) createICMPEchoRequest() []byte {
    body := &icmp.Echo{
        ID:   p.ID,
        Seq:  p.Seq,
        Data: make([]byte, p.Options.Size),
    }
    message := &icmp.Message{
        Type: ipv4.ICMPTypeEcho,
        Code: 0,
        Body: body,
    }
    msgBytes, err := message.Marshal(nil)
    if err != nil {
        logger.Logger.Fatalf("Failed to marshal ICMP message: %v", err)
    }

    // TODO: Implement Record Route and Timestamp options if enabled
    // This requires constructing the IP options manually, which is beyond the scope of this example.

    return msgBytes
}

func (p *Pinger) getTTL(peer net.Addr) int {
    // Placeholder for retrieving TTL from the reply
    // In Go's net package, retrieving TTL is not straightforward and may require low-level socket options
    return 64 // Default TTL value
}

func (p *Pinger) Statistics() {
    p.Stats.Lost = p.Stats.Transmitted - p.Stats.Received
    loss := float64(p.Stats.Lost) / float64(p.Stats.Transmitted) * 100
    var sum time.Duration
    for _, rtt := range p.RTTs {
        sum += rtt
        if rtt < p.Stats.MinRTT {
            p.Stats.MinRTT = rtt
        }
        if rtt > p.Stats.MaxRTT {
            p.Stats.MaxRTT = rtt
        }
    }
    if p.Stats.Received > 0 {
        p.Stats.AvgRTT = sum / time.Duration(p.Stats.Received)
    }

    fmt.Printf("\n--- %s ping statistics ---\n", p.Destination)
    fmt.Printf("%d packets transmitted, %d packets received, %.0f%% packet loss\n",
        p.Stats.Transmitted, p.Stats.Received, loss)
    if p.Stats.Received > 0 {
        fmt.Printf("round-trip min/avg/max = %.3f/%.3f/%.3f ms\n",
            float64(p.Stats.MinRTT.Microseconds())/1000.0,
            float64(p.Stats.AvgRTT.Microseconds())/1000.0,
            float64(p.Stats.MaxRTT.Microseconds())/1000.0)
    }
}