// snmp_manager/table_handler.go

package snmp_manager

import (
    "log"

    "github.com/gosnmp/gosnmp"
)

// Table represents an SNMP table with multiple columns.
type Table struct {
    Name     string
    Columns  []string
    Variables map[string]MIBVariable
}

// TableManager manages multiple SNMP tables.
type TableManager struct {
    Tables map[string]Table
}

// NewTableManager initializes a new TableManager.
func NewTableManager() *TableManager {
    return &TableManager{
        Tables: make(map[string]Table),
    }
}

// AddTable adds a table to the manager.
func (tm *TableManager) AddTable(table Table) {
    tm.Tables[table.Name] = table
    Logger.Printf("Table %s added to TableManager", table.Name)
}

// GetTable returns the table with the given name.
func (tm *TableManager) GetTable(name string) (Table, bool) {
    table, exists := tm.Tables[name]
    return table, exists
}

// Example usage: Define and add tables based on parsed MIBs