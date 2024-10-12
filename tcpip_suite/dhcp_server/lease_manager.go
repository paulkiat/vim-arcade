// Manages IP leases, including assignment, renewal, and expiration.

// dhcp_server/lease_manager.go

package lease_manager

import (
    "log"
    "net"
    "strings"
    "sync"
    "time"
)

type Lease struct {
    IP        string
    MAC       string
    Expiry    time.Time
}

type LeaseManager struct {
    ipPool     []string
    leases     map[string]Lease // key: MAC
    leaseTime  time.Duration
    mutex      sync.Mutex
    logger     *log.Logger
}

func NewLeaseManager(startIP, endIP string, leaseTime int, logger *log.Logger) *LeaseManager {
    manager := &LeaseManager{
        ipPool:    generateIPPool(startIP, endIP),
        leases:    make(map[string]Lease),
        leaseTime: time.Duration(leaseTime) * time.Second,
        logger:    logger,
    }

    go manager.cleanupLeases()
    return manager
}

func generateIPPool(startIP, endIP string) []string {
    ips := []string{}
    start := net.ParseIP(startIP).To4()
    end := net.ParseIP(endIP).To4()

    for ip := start; !ip.Equal(end.Add(1)); incrementIP(ip) {
        ips = append(ips, ip.String())
    }
    return ips
}

func incrementIP(ip net.IP) {
    for j := len(ip) - 1; j >= 0; j-- {
        ip[j]++
        if ip[j] > 0 {
            break
        }
    }
}

func (lm *LeaseManager) AssignIP(mac string) string {
    lm.mutex.Lock()
    defer lm.mutex.Unlock()

    mac = strings.ToLower(mac)
    if lease, exists := lm.leases[mac]; exists {
        if time.Now().Before(lease.Expiry) {
            lm.logger.Printf("Renewing lease for MAC %s: IP %s", mac, lease.IP)
            lm.leases[mac] = Lease{
                IP:     lease.IP,
                MAC:    mac,
                Expiry: time.Now().Add(lm.leaseTime),
            }
            return lease.IP
        }
    }

    if len(lm.ipPool) == 0 {
        lm.logger.Println("IP pool exhausted.")
        return ""
    }

    assignedIP := lm.ipPool[0]
    lm.ipPool = lm.ipPool[1:]

    lm.leases[mac] = Lease{
        IP:     assignedIP,
        MAC:    mac,
        Expiry: time.Now().Add(lm.leaseTime),
    }

    lm.logger.Printf("Assigned IP %s to MAC %s", assignedIP, mac)
    return assignedIP
}

func (lm *LeaseManager) GetIP(mac string) string {
    lm.mutex.Lock()
    defer lm.mutex.Unlock()

    mac = strings.ToLower(mac)
    if lease, exists := lm.leases[mac]; exists {
        if time.Now().Before(lease.Expiry) {
            return lease.IP
        }
    }
    return ""
}

func (lm *LeaseManager) ReleaseIP(mac string) {
    lm.mutex.Lock()
    defer lm.mutex.Unlock()

    mac = strings.ToLower(mac)
    if lease, exists := lm.leases[mac]; exists {
        lm.ipPool = append(lm.ipPool, lease.IP)
        delete(lm.leases, mac)
        lm.logger.Printf("Released IP %s from MAC %s", lease.IP, mac)
    }
}

func (lm *LeaseManager) cleanupLeases() {
    for {
        time.Sleep(60 * time.Second)
        lm.mutex.Lock()
        now := time.Now()
        for mac, lease := range lm.leases {
            if now.After(lease.Expiry) {
                lm.logger.Printf("Lease expired: IP %s from MAC %s", lease.IP, mac)
                lm.ipPool = append(lm.ipPool, lease.IP)
                delete(lm.leases, mac)
            }
        }
        lm.mutex.Unlock()
    }
}

func (lm *LeaseManager) Stop() {
    // Placeholder for any cleanup if necessary
}
