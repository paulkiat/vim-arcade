// Manages loaded MIBs and provides utility functions to interact with MIB variables.
// snmp_manager/mib_manager.go

package snmp_manager

import (
    "fmt"
    "strings"
)

// MIBManager manages multiple MIBs and provides lookup functionalities.
type MIBManager struct {
    MIBs map[string]MIB
}

// NewMIBManager initializes a new MIBManager.
func NewMIBManager() *MIBManager {
    return &MIBManager{
        MIBs: make(map[string]MIB),
    }
}

// AddMIB adds a parsed MIB to the manager.
func (m *MIBManager) AddMIB(mib MIB) {
    m.MIBs[mib.Name] = mib
    Logger.Printf("MIB %s added to MIBManager", mib.Name)
}

// GetObjectID returns the Object Identifier for a given MIB variable.
func (m *MIBManager) GetObjectID(mibName string, varName string) (string, error) {
    mib, exists := m.MIBs[mibName]
    if !exists {
        return "", fmt.Errorf("MIB %s not found", mibName)
    }

    for _, variable := range mib.Variables {
        if strings.EqualFold(variable.Name, varName) {
            return variable.ObjectID, nil
        }
    }

    return "", fmt.Errorf("variable %s not found in MIB %s", varName, mibName)
}

// Example usage: GetObjectID("MIB-II", "udpInDatagrams")