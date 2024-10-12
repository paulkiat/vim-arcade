// snmp_manager/trap_handler.go

package snmp_manager

import (
    "log"

    "github.com/gosnmp/gosnmp"
)

// HandleTraps initializes a trap listener to receive SNMP traps.
func (m *Manager) HandleTraps() {
    defer m.wg.Done()

    trapListener := &gosnmp.TrapListener{
        OnNewTrap: func(packet *gosnmp.SnmpPacket, addr *gosnmp.SnmpAddr) {
            Logger.Printf("Received trap from %s", addr.IP)
            Logger.Printf("Community: %s", packet.Community)
            Logger.Printf("Enterprise: %s", packet.Enterprise)
            Logger.Printf("Agent Address: %s", packet.AgentAddress)
            Logger.Printf("Generic Trap: %d", packet.GenericTrap)
            Logger.Printf("Specific Trap: %d", packet.SpecificTrap)
            Logger.Printf("Timestamp: %d", packet.TimeTicks)

            for _, variable := range packet.Variables {
                switch variable.Type {
                case gosnmp.OctetString:
                    Logger.Printf("OID: %s, Value: %s", variable.Name, string(variable.Value.([]byte)))
                case gosnmp.Integer:
                    Logger.Printf("OID: %s, Value: %d", variable.Name, gosnmp.ToBigInt(variable.Value).Int64())
                case gosnmp.Counter32, gosnmp.Counter64, gosnmp.Gauge32, gosnmp.TimeTicks, gosnmp.ObjectIdentifier:
                    Logger.Printf("OID: %s, Value: %v", variable.Name, variable.Value)
                default:
                    Logger.Printf("OID: %s, Value: %v", variable.Name, variable.Value)
                }
            }
        },
    }

    // Start listening for traps on UDP port 162
    err := trapListener.Listen("0.0.0.0:162")
    if err != nil {
        Logger.Fatalf("Failed to listen for traps: %v", err)
    }

    Logger.Println("Trap listener started on port 162.")
    select {} // Block indefinitely
}

// SendTrap allows sending custom traps to SNMP agents (optional feature).
func (m *Manager) SendTrap(agent *gosnmp.GoSNMP, trapType int, variables []gosnmp.SnmpPDU) error {
    trap := gosnmp.SnmpPacket{
        Variables:       variables,
        Enterprise:      "1.3.6.1.4.1.12345", // Example enterprise OID
        AgentAddress:    agent.Target,
        GenericTrap:     trapType,
        SpecificTrap:    0,
        TimeTicks:       uint32(time.Now().Unix()),
        Community:       m.config.Community,
        Version:         gosnmp.Version1,
        SecurityModel:   gosnmp.NoAuthNoPriv,
        SecurityParameters: &gosnmp.UsmSecurityParameters{},
    }

    _, err := agent.SendTrap(trapType, variables)
    if err != nil {
        return err
    }

    Logger.Printf("Sent trap type %d to agent %s:%d", trapType, agent.Target, agent.Port)
    return nil
}