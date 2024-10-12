// pkg/ethernet/ethernet.go

package ethernet

import (
    "encoding/binary"
    "errors"
    "hash/crc32"
    "net"

    "link_layer/pkg/utils"
)

// EthernetFrame represents an Ethernet frame.
type EthernetFrame struct {
    DestinationMAC net.HardwareAddr
    SourceMAC      net.HardwareAddr
		VLANID         uint16 // 0 if no VLAN tagging is used.
    EtherType      uint16
    Payload        []byte
    CRC            uint32
}

// Constants for EtherType values.
const (
    EtherTypeIPv4   = 0x0800
    EtherTypeARP    = 0x0806
    EtherTypeIPv6   = 0x86DD
    EtherTypeMACSec = 0x88E5
    EtherType8021Q  = 0x8100 // VLAN tagging
)

// WorkerPool manages a pool of worker goroutines.
type WorkerPool struct {
	workers int
	jobs    chan *EthernetFrame
	wg      sync.WaitGroup
}

// NewWorkerPool initializes a new WorkerPool.
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
			workers: workers,
			jobs:    make(chan *EthernetFrame, 100),
	}
}

// Start launches the worker pool.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
			wp.wg.Add(1)
			go func(workerID int) {
					defer wp.wg.Done()
					for frame := range wp.jobs {
							serialized, err := frame.Serialize()
							if err != nil {
									utils.Logger.Errorf("Worker %d: Failed to serialize frame: %v", workerID, err)
									continue
							}
							// Handle the serialized frame (e.g., send over network)
							utils.Logger.Infof("Worker %d: Serialized Ethernet Frame: %x", workerID, serialized)
					}
			}(i)
	}
}

// AddJob adds a frame to the job queue.
func (wp *WorkerPool) AddJob(frame *EthernetFrame) {
	wp.jobs <- frame
}

// Stop gracefully shuts down the worker pool.
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}

// NewEthernetFrame creates a new Ethernet frame.
func NewEthernetFrame(destMAC, srcMAC net.HardwareAddr, etherType uint16, payload []byte, vlanID uint16 ...uint16) *EthernetFrame {
    return &EthernetFrame{
        DestinationMAC: destMAC,
        SourceMAC:      srcMAC,
				VLANID:         vlanID,
        EtherType:      etherType,
        Payload:        payload,
        CRC:            0, // Placeholder for CRC calculation.
    }
}

// Serialize converts the EthernetFrame into bytes.
func (ef *EthernetFrame) Serialize() ([]byte, error) {
    if len(ef.DestinationMAC) != 6 || len(ef.SourceMAC) != 6 {
        return nil, errors.New("invalid MAC address length")
    }

		var frame []byte
		if ef.VLANID != 0 { // VLAN tagging
			frame = make([]byte, 18 + len(ef.Payload) + 4) // 14 bytes header, 20 bytes IPv4 header, 4 bytes CRC.
			copy(frame[0:6], ef.DestinationMAC)
			copy(frame[6:12], ef.SourceMAC)
			binary.BigEndian.PutUint16(frame[12:14], EtherTypeDot1Q)
			binary.BigEndian.PutUint16(frame[14:16], ef.VLANID)
			binary.BigEndian.PutUint16(frame[16:18], ef.EtherType)
			copy(frame[18:18+len(ef.Payload)], ef.Payload)
		} else {
      // Standard Ethernet frame: 14 bytes header
			frame = make([]byte, 14+len(ef.Payload)+4)
			copy(frame[0:6], ef.DestinationMAC)
			copy(frame[6:12], ef.SourceMAC)
			binary.BigEndian.PutUint16(frame[12:14], ef.EtherType)
			copy(frame[14:14+len(ef.Payload)], ef.Payload)
		}

		if frame == nil {
			return nil, errors.New("invalid EtherType");
		}

		// Ethernet frame
		// 0-5: Destination MAC
		// 6-11: Source MAC
		// 12-13: EtherType
		// 14: Payload
		// 15-18: CRC
		if ef.EtherType == EtherTypeIPv4 {
			frame = make([]byte, 14+20+len(ef.Payload)+4) // 14 bytes header, 20 bytes IPv4 header, payload, 4 bytes CRC.
		} else if ef.EtherType == EtherTypeARP {
			frame = make([]byte, 14+28+len(ef.Payload)+4) // 14 bytes header, 28 bytes ARP header, payload, 4 bytes CRC.
		}

		// todo: add IPv6 ??

		if frame == nil {
				return nil, errors.New("invalid EtherType")
		}

    frame := make([]byte, 14+len(ef.Payload)+4) // 14 bytes header, payload, 4 bytes CRC.

    copy(frame[0:6], ef.DestinationMAC)
    copy(frame[6:12], ef.SourceMAC)
    binary.BigEndian.PutUint16(frame[12:14], ef.EtherType)
    copy(frame[14:14+len(ef.Payload)], ef.Payload)

    // CRC calculation using CRC32.
    ef.CRC = crc32.ChecksumIEEE(frame[:14+len(ef.Payload)])
    binary.BigEndian.PutUint32(frame[14+len(ef.Payload):], ef.CRC)

    return frame, nil
}

// DeserializeEthernetFrame parses bytes into an EthernetFrame.
func DeserializeEthernetFrame(data []byte) (*EthernetFrame, error) {
    if len(data) < 18 { // Minimum 14 bytes header + 4 bytes CRC.
        return nil, errors.New("data too short to be a valid Ethernet frame")
    }

    frame := &EthernetFrame{
        DestinationMAC: net.HardwareAddr(data[0:6]),
        SourceMAC:      net.HardwareAddr(data[6:12]),
        EtherType:      binary.BigEndian.Uint16(data[12:14]),
        Payload:        data[14 : len(data)-4],
        CRC:            binary.BigEndian.Uint32(data[len(data)-4:]),
    }

    // Verify CRC
    calculatedCRC := crc32.ChecksumIEEE(data[:len(data)-4])
    if calculatedCRC != frame.CRC {
        return nil, errors.New("CRC mismatch")
    }

    return frame, nil
}

// CalculateCRC calculates the CRC32 for the given data.
func CalculateCRC(data []byte) uint32 {
	return crc32Checksum(data)
}

// crc32Checksum computes the CRC32 checksum using IEEE polynomial.
func crc32Checksum(data []byte) uint32 {
	table := crc32.MakeTable(crc32.IEEE)
	return crc32.Checksum(data, table)
}

// AddVLANTag adds a VLAN tag to the Ethernet frame.
func (ef *EthernetFrame) AddVLANTag(vlanID uint16) error {
	if ef.EtherType != EtherTypeIPv4 && ef.EtherType != EtherTypeIPv6 && ef.EtherType != EtherTypeARP {
			return errors.New("unsupported EtherType for VLAN tagging")
	}

	// Insert VLAN tag after SourceMAC and before EtherType.
	vlanTag := make([]byte, 4)
	binary.BigEndian.PutUint16(vlanTag[0:2], EtherType8021Q)
	binary.BigEndian.PutUint16(vlanTag[2:4], vlanID)

	// Reconstruct frame with VLAN tag.
	newFrame := make([]byte, 14+4+len(ef.Payload)+4)
	copy(newFrame[0:6], ef.DestinationMAC)
	copy(newFrame[6:12], ef.SourceMAC)
	copy(newFrame[12:14], vlanTag[0:2])
	copy(newFrame[14:18], vlanTag[2:4])
	copy(newFrame[18:18+len(ef.Payload)], ef.Payload)

	// Update EtherType.
	ef.EtherType = EtherType8021Q

	// Recalculate CRC.
	ef.CRC = CalculateCRC(newFrame[:18+len(ef.Payload)])
	binary.BigEndian.PutUint32(newFrame[18+len(ef.Payload):], ef.CRC)

	// Update Payload to include VLAN tag.
	ef.Payload = newFrame[18 : 18+len(ef.Payload)]

	return nil
}
