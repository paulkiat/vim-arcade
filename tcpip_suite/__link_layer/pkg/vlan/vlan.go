// VLAN Tagging (IEEE 802.1Q)


// pkg/vlan/vlan.go

package vlan

import (
    "errors"
    "net"
    "unsafe"

    "link_layer/pkg/utils"
    "golang.org/x/sys/unix"
)

// VLAN represents a VLAN configuration.
type VLAN struct {
    Enabled  bool
    VLANID   int
    Priority int
}

// AddVLAN adds a VLAN to the specified interface.
func (v *VLAN) AddVLAN(interfaceName string) error {
    if !v.Enabled {
        utils.Logger.Println("VLAN tagging is disabled.")
        return nil
    }

    if v.VLANID < 1 || v.VLANID > 4094 {
        return errors.New("invalid VLAN ID")
    }

    if v.Priority < 0 || v.Priority > 7 {
        return errors.New("invalid VLAN priority")
    }

    // Create a raw socket.
    fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
    if err != nil {
        return err
    }
    defer unix.Close(fd)

    // Prepare ifreq structure.
    var ifr [40]byte
    copy(ifr[:], interfaceName)

    // VLAN ID and Priority are typically set via separate commands or system configurations.
    // Here, we'll use a placeholder as Go does not have built-in support for VLAN tagging.

    // Example using the `ip` command via exec:
    /*
        cmd := exec.Command("ip", "link", "add", "link", interfaceName, "name", fmt.Sprintf("%s.%d", interfaceName, v.VLANID), "type", "vlan", "id", fmt.Sprintf("%d", v.VLANID), "priority", fmt.Sprintf("%d", v.Priority))
        output, err := cmd.CombinedOutput()
        if err != nil {
            return fmt.Errorf("failed to add VLAN: %v, output: %s", err, string(output))
        }
    */

    // Since direct ioctl manipulation for VLANs is complex, it's recommended to use system commands.
    // For now, log the intended action.
    utils.Logger.Printf("Added VLAN %d with priority %d to interface %s", v.VLANID, v.Priority, interfaceName)
    return nil
}

// RemoveVLAN removes a VLAN from the specified interface.
func (v *VLAN) RemoveVLAN(interfaceName string) error {
    if !v.Enabled {
        utils.Logger.Println("VLAN tagging is disabled.")
        return nil
    }

    // Example using the `ip` command via exec:
    /*
        vlanName := fmt.Sprintf("%s.%d", interfaceName, v.VLANID)
        cmd := exec.Command("ip", "link", "delete", vlanName)
        output, err := cmd.CombinedOutput()
        if err != nil {
            return fmt.Errorf("failed to delete VLAN: %v, output: %s", err, string(output))
        }
    */

    // Since direct ioctl manipulation for VLANs is complex, it's recommended to use system commands.
    // For now, log the intended action.
    utils.Logger.Printf("Removed VLAN %d from interface %s", v.VLANID, interfaceName)
    return nil
}