// pkg/mtu/mtu.go

package mtu

import (
    "errors"
    "net"
    "unsafe"

    "link_layer/pkg/utils"
    "golang.org/x/sys/unix"
)

// GetMTU retrieves the MTU of a given network interface.
func GetMTU(ifaceName string) (int, error) {
    iface, err := net.InterfaceByName(ifaceName)
    if err != nil {
        return 0, err
    }
    return iface.MTU, nil
}

// SetMTU sets the MTU of a given network interface.
// This implementation is Linux-specific.
func SetMTU(ifaceName string, mtu int) error {
    if mtu < 68 || mtu > 9000 {
        return errors.New("invalid MTU size")
    }

    // Create a raw socket.
    fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
    if err != nil {
        return err
    }
    defer unix.Close(fd)

    // Prepare ifreq structure.
    var ifr [40]byte
    copy(ifr[:], ifaceName)
    binary.BigEndian.PutUint32(ifr[16:], uint32(mtu))

    // Perform ioctl to set MTU.
    _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCSIFMTU), uintptr(unsafe.Pointer(&ifr[0])))
    if errno != 0 {
        return errno
    }

    utils.Logger.Printf("Set MTU of interface %s to %d", ifaceName, mtu)
    return nil
}