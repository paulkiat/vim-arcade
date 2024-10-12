// snmp_manager/snmp.go

package snmp_manager

import (
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/gosnmp/gosnmp"
    "gopkg.in/yaml.v2"
)

// Config holds the SNMP Manager configuration.
type Config struct {
    Version string        `yaml:"version"`
    Community string      `yaml:"community"`
    Users   []SNMPUser    `yaml:"users"`
    Agents  []SNMPAgent   `yaml:"agents"`
    MIBs    []MIB         `yaml:"mibs"`
    Logging LoggingConfig `yaml:"logging"`
}

// SNMPUser defines the SNMPv3 user credentials.
type SNMPUser struct {
    Username       string `yaml:"username"`
    AuthProtocol   string `yaml:"auth_protocol"` // e.g., "SHA", "MD5"
    AuthPassword   string `yaml:"auth_password"`
    PrivProtocol   string `yaml:"priv_protocol"` // e.g., "AES", "DES"
    PrivPassword   string `yaml:"priv_password"`
}

// SNMPAgent defines the SNMP agent details.
type SNMPAgent struct {
    IP   string `yaml:"ip"`
    Port uint16 `yaml:"port"`
}

// MIB represents a Management Information Base.
type MIB struct {
    Name string `yaml:"name"`
    File string `yaml:"file"`
}

// LoggingConfig holds logging-related configurations.
type LoggingConfig struct {
    Level string `yaml:"level"`
    File  string `yaml:"file"`
}

// Manager represents the SNMP Manager.
type Manager struct {
		config     Config
		clients    []*gosnmp.GoSNMP
		wg         sync.WaitGroup
		quit       chan bool
		mibManager *MIBManager
		tableManager *TableManager
}

// NewManager initializes the SNMP Manager with the given configuration.
func NewManager(config Config) (*Manager, error) {
    manager := &Manager{
				config:     config,
				clients:    []*gosnmp.GoSNMP{},
				quit:       make(chan bool),
				mibManager: NewMIBManager(),
				tableManager: NewTableManager()
    }

    // Initialize SNMP clients based on version
    for _, agent := range config.Agents {
        if config.Version == "1" || config.Version == "2c" {
            g := &gosnmp.GoSNMP{
                Target:    agent.IP,
                Port:      agent.Port,
                Version:   gosnmp.Version1,
                Community: config.Community,
                Timeout:   time.Duration(5) * time.Second,
                Retries:   3,
                Logger:    log.Default(),
            }

            err := g.Connect()
            if err != nil {
                Logger.Printf("Error connecting to agent %s: %v", agent.IP, err)
                continue
            }

            manager.clients = append(manager.clients, g)
            Logger.Printf("Connected to SNMP agent %s:%d", agent.IP, agent.Port)
        } else if config.Version == "3" {
            for _, user := range config.Users {
                g := &gosnmp.GoSNMP{
                    Target:             agent.IP,
                    Port:               agent.Port,
                    Version:            gosnmp.Version3,
                    Timeout:            time.Duration(5) * time.Second,
                    Retries:            3,
                    SecurityModel:      gosnmp.UserSecurityModel,
                    MsgFlags:           gosnmp.AuthPriv, // Authentication and Privacy
                    SecurityParameters: &gosnmp.UsmSecurityParameters{},
                    Logger:             log.Default(),
                }

                // Set SNMPv3 security parameters
                g.SecurityParameters = &gosnmp.UsmSecurityParameters{
                    UserName:                 user.Username,
                    AuthenticationProtocol:   gosnmp.StringToSNMPv3AuthProtocol(user.AuthProtocol),
                    AuthenticationPassphrase: user.AuthPassword,
                    PrivacyProtocol:          gosnmp.StringToSNMPv3PrivProtocol(user.PrivProtocol),
                    PrivacyPassphrase:        user.PrivPassword,
                }

                err := g.Connect()
                if err != nil {
                    Logger.Printf("Error connecting to agent %s: %v", agent.IP, err)
                    continue
                }

                manager.clients = append(manager.clients, g)
                Logger.Printf("Connected to SNMP agent %s:%d as user %s", agent.IP, agent.Port, user.Username)
            }
        }
    }

    if len(manager.clients) == 0 {
        return nil, fmt.Errorf("no SNMP agents connected successfully")
    }

    // Load MIBs
    for _, mib := range config.MIBs {
			parsedMIB, err := parseMIB(mib.File)
			if err != nil {
					Logger.Printf("Failed to parse MIB %s: %v", mib.Name, err)
					continue
			}
			manager.mibManager.AddMIB(parsedMIB)
	}


    // Initialize tables based on MIB variables
    for _, mib := range manager.mibManager.MIBs {
			for _, variable := range mib.Variables {
					if strings.HasSuffix(variable.Name, "Table") {
							// Identify table columns (assuming variables with same prefix)
							tableName := variable.Name
							table := Table{
									Name:      tableName,
									Columns:   []string{},
									Variables: make(map[string]MIBVariable),
							}

							// Iterate through MIB variables to find columns belonging to this table
							for _, varb := range mib.Variables {
									if strings.HasPrefix(varb.Name, tableName) && !strings.Contains(varb.Name, "Table") {
											table.Columns = append(table.Columns, varb.Name)
											table.Variables[varb.Name] = varb
									}
							}

							if len(table.Columns) > 0 {
									manager.tableManager.AddTable(table)
							}
					}
			}
	}
    return manager, nil
}

// loadMIB loads and parses a MIB file. (Placeholder for MIB parsing implementation)
func loadMIB(filePath string) error {
    // Implement MIB parsing logic or integrate with existing MIB parsing libraries
    // For simplicity, this function is a placeholder
    Logger.Printf("Loading MIB file: %s", filePath)
    return nil
}

// Run starts the SNMP Manager operations.
func (m *Manager) Run() {
    Logger.Println("SNMP Manager is running.")

    // Start trap listener if SNMPv1/v2c does not support traps, else handled separately
    m.wg.Add(1)
    go m.HandleTraps()

    // Periodically fetch SNMP variables
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()

    m.wg.Add(1)
    go func() {
        defer m.wg.Done()
        for {
            select {
            case <-ticker.C:
                m.FetchVariables()
            case <-m.quit:
                return
            }
        }
    }()
}

// FetchVariables performs SNMP get operations for predefined OIDs.
func (m *Manager) FetchVariables() {
    // Define OIDs based on MIB groups. Example for udpInDatagrams, udpNoPorts, etc.
    oids := []string{
        "1.3.6.1.2.1.7.1.0", // udpInDatagrams.0
        "1.3.6.1.2.1.7.2.0", // udpNoPorts.0
        "1.3.6.1.2.1.7.3.0", // udpInErrors.0
        "1.3.6.1.2.1.7.4.0", // udpOutDatagrams.0
        // Add more OIDs as needed for other MIB groups
    }

    for _, client := range m.clients {
        go func(c *gosnmp.GoSNMP) {
            Logger.Printf("Fetching variables from agent %s:%d", c.Target, c.Port)
            result, err := c.Get(oids) // Get() accepts up to 8 OIDs
            if err != nil {
                Logger.Printf("SNMP Get failed for %s: %v", c.Target, err)
                return
            }

            for _, variable := range result.Variables {
                switch variable.Type {
                case gosnmp.OctetString:
                    Logger.Printf("%s = %s", variable.Name, string(variable.Value.([]byte)))
                case gosnmp.Integer:
                    Logger.Printf("%s = %d", variable.Name, gosnmp.ToBigInt(variable.Value).Int64())
                case gosnmp.Counter32, gosnmp.Counter64, gosnmp.Gauge32, gosnmp.TimeTicks:
                    Logger.Printf("%s = %v", variable.Name, variable.Value)
                default:
                    Logger.Printf("%s = %v", variable.Name, variable.Value)
                }
            }
        }(client)
    }
}

// GetBulkVariables performs SNMP get-bulk operations for efficient data retrieval.
func (m *Manager) GetBulkVariables(oids []string, maxRepetitions int) {
    for _, client := range m.clients {
        go func(c *gosnmp.GoSNMP) {
            Logger.Printf("Performing GetBulk from agent %s:%d", c.Target, c.Port)
            result, err := c.GetBulk(oids, maxRepetitions, 0)
            if err != nil {
                Logger.Printf("SNMP GetBulk failed for %s: %v", c.Target, err)
                return
            }

            for _, variable := range result.Variables {
                switch variable.Type {
                case gosnmp.OctetString:
                    Logger.Printf("OID: %s, Value: %s", variable.Name, string(variable.Value.([]byte)))
                case gosnmp.Integer:
                    Logger.Printf("OID: %s, Value: %d", variable.Name, gosnmp.ToBigInt(variable.Value).Int64())
                case gosnmp.Counter32, gosnmp.Counter64, gosnmp.Gauge32, gosnmp.TimeTicks:
                    Logger.Printf("OID: %s, Value: %v", variable.Name, variable.Value)
                default:
                    Logger.Printf("OID: %s, Value: %v", variable.Name, variable.Value)
                }
            }
        }(client)
    }
}

// SetVariable sets the value of a specific SNMP variable.
func (m *Manager) SetVariable(oid string, value interface{}) {
    for _, client := range m.clients {
        go func(c *gosnmp.GoSNMP) {
            Logger.Printf("Setting OID %s on agent %s:%d", oid, c.Target, c.Port)
            pdu := gosnmp.SnmpPDU{
                Name:  oid,
                Type:  gosnmp.Integer, // Adjust type based on variable
                Value: value,
            }

            result, err := c.Set([]gosnmp.SnmpPDU{pdu})
            if err != nil {
                Logger.Printf("SNMP Set failed for %s: %v", c.Target, err)
                return
            }

            for _, variable := range result.Variables {
                Logger.Printf("Set OID: %s, Value: %v", variable.Name, variable.Value)
            }
        }(client)
    }
}

// Stop gracefully stops the SNMP Manager.
func (m *Manager) Stop() {
    close(m.quit)
    m.wg.Wait()

    for _, client := range m.clients {
        client.Conn.Close()
    }

    Logger.Println("SNMP Manager has stopped.")
}