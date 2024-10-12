// pkg/arp/arp.go

package arp

import (
    "encoding/binary"
    "errors"
    "net"

    "link_layer/pkg/utils"
)

// ARPPacket represents an ARP packet.
type ARPPacket struct {
    HardwareType uint16
    ProtocolType uint16
    HardwareSize uint8
    ProtocolSize uint8
    Operation    uint16
    SenderMAC    net.HardwareAddr
    SenderIP     net.IP
    TargetMAC    net.HardwareAddr
    TargetIP     net.IP
}

// Constants for ARP operations.
const (
    ARPRequest = 1
    ARPReply   = 2
)

// NewARPPacket creates a new ARP request or reply.
func NewARPPacket(op uint16, senderMAC net.HardwareAddr, senderIP, targetIP net.IP, targetMAC net.HardwareAddr) *ARPPacket {
    return &ARPPacket{
        HardwareType: 1,         // Ethernet
        ProtocolType: 0x0800,    // IPv4
        HardwareSize: 6,
        ProtocolSize: 4,
        Operation:    op,
        SenderMAC:    senderMAC,
        SenderIP:     senderIP,
        TargetMAC:    targetMAC,
        TargetIP:     targetIP,
    }
}

// Serialize converts the ARPPacket into bytes.
func (arp *ARPPacket) Serialize() ([]byte, error) {
    if len(arp.SenderMAC) != 6 || len(arp.TargetMAC) != 6 {
        return nil, errors.New("invalid MAC address length")
    }
    if len(arp.SenderIP) != 4 || len(arp.TargetIP) != 4 {
        return nil, errors.New("invalid IP address length")
    }

    packet := make([]byte, 28)
    binary.BigEndian.PutUint16(packet[0:2], arp.HardwareType)
    binary.BigEndian.PutUint16(packet[2:4], arp.ProtocolType)
    packet[4] = arp.HardwareSize
    packet[5] = arp.ProtocolSize
    binary.BigEndian.PutUint16(packet[6:8], arp.Operation)
    copy(packet[8:14], arp.SenderMAC)
    copy(packet[14:18], arp.SenderIP.To4())
    copy(packet[18:24], arp.TargetMAC)
    copy(packet[24:28], arp.TargetIP.To4())

    return packet, nil
}

// DeserializeARPPacket parses bytes into an ARPPacket.
func DeserializeARPPacket(data []byte) (*ARPPacket, error) {
    if len(data) < 28 {
        return nil, errors.New("data too short to be a valid ARP packet")
    }

    packet := &ARPPacket{
        HardwareType: binary.BigEndian.Uint16(data[0:2]),
        ProtocolType: binary.BigEndian.Uint16(data[2:4]),
        HardwareSize: data[4],
        ProtocolSize: data[5],
        Operation:    binary.BigEndian.Uint16(data[6:8]),
        SenderMAC:    net.HardwareAddr(data[8:14]),
        SenderIP:     net.IP(data[14:18]),
        TargetMAC:    net.HardwareAddr(data[18:24]),
        TargetIP:     net.IP(data[24:28]),
    }

    return packet, nil
}