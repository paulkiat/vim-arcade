package packet

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"vim-arcade.theprimeagen.com/pkg/assert"
	"vim-arcade.theprimeagen.com/pkg/utils"
)

const VERSION = 1

const HEADER_SIZE = 4
const TYPE_ENC_INDEX = 1
const HEADER_LENGTH_OFFSET = 2
const PACKET_MAX_SIZE = 1024
const PACKET_PAYLOAD_SIZE = 1024 - HEADER_SIZE

var PacketMaxSizeExceeded = fmt.Errorf("Packet length has exceeded allowed size of %d", PACKET_PAYLOAD_SIZE - 1)
var PacketVersionMismatch = fmt.Errorf("Expected packet version to equal %d", VERSION)
var PacketBufferNotBigEnough = fmt.Errorf("Buffer could not fit the entire packet")

type Encoding uint8

const (
    EncodingJSON Encoding = iota
    EncodingCustom
    EncodingUNUSED
    EncodingUNUSED2
)

type Packet struct {
    data []byte
    len int
}

type PacketEncoder interface {
    io.Reader
    Type() uint8
    Encoding() Encoding
}

func ValidHeader(data []byte) bool {
    return data[0] == VERSION
}

func ContainsCompletePacket(data []byte) bool {
    if len(data) >= HEADER_SIZE {
        return false
    }

    l := getPacketLength(data)
    return len(data) >= int(l) + HEADER_SIZE
}

func getPacketLength(data []byte) uint16 {
    return binary.BigEndian.Uint16(data[HEADER_LENGTH_OFFSET:])
}

func PacketFromBytes(data []byte) Packet {
    assert.Assert(data[0] == VERSION, "version mismatch: this should be handled by the framer before packet is created", "VERSION", VERSION, "provided", data[0])

    // TODO will there ever be a 0 length buffer
    dataLen := len(data) - HEADER_SIZE
    assert.Assert(dataLen > 0, "packets must contain some sort of data")

    encodedLen := getPacketLength(data)
    assert.Assert(dataLen == int(encodedLen), "the data buffer provided has a length mismatch", "expected length", dataLen, "encoded length", encodedLen)

    return Packet{
        data: data,
        len: len(data),
    }
}

func NewPacket(encoder PacketEncoder) Packet {
    b := make([]byte, PACKET_MAX_SIZE, PACKET_MAX_SIZE)

    enc := b[HEADER_SIZE:]
    n, err := encoder.Read(enc)
    assert.NoError(err, "i should never fail on encoding a packet")
    assert.Assert(n != PACKET_PAYLOAD_SIZE, "max packet size exceeded", "MAX_SIZE", PACKET_PAYLOAD_SIZE)

    t := encoder.Type()
    assert.Assert(t <= (0x40 - 1), "type has exceeded allowed size", "type", t)

    b[0] = VERSION
    b[1] = uint8(encoder.Encoding() << 6) | encoder.Type()

    binary.BigEndian.PutUint16(b[2:], uint16(n))

    return Packet{data: b, len: n + HEADER_SIZE}
}

func (p *Packet) Into(writer io.Writer) (int, error) {
    // TODO v2
    // reconsider just trusting the science and if it fails, assert and report
    // to golang's github
    return p.len, utils.WriteAll(p.data[:p.len], writer)
}

func (p *Packet) Data() []byte {
    return p.data[HEADER_SIZE:p.len]
}

func (p *Packet) Type() uint8 {
    return p.data[TYPE_ENC_INDEX] & 0x3F
}

func (p *Packet) Encoding() Encoding {
    return Encoding((p.data[TYPE_ENC_INDEX] >> 6) & 0x3)
}

func (p *Packet) Read(data []byte) (int, error) {
    if len(data) < p.len {
        return 0, PacketBufferNotBigEnough
    }
    copy(data, p.data[0:p.len])
    return p.len, nil
}

func (p *Packet) String() string {
    return fmt.Sprintf("Packet(%d) -> \"%s\"", p.len, string(p.data))
}

type PacketFramer struct {
    buf []byte
    idx int
}

func NewPacketFramer() PacketFramer {
    return PacketFramer{
        buf: make([]byte, PACKET_PAYLOAD_SIZE, PACKET_PAYLOAD_SIZE),
    }
}

func (p *PacketFramer) Push(data []byte) {
    n := copy(p.buf[p.idx:], data)

    if n < len(data) {
        p.buf = append(p.buf, data[n:]...)
    }
    p.idx += len(data)
}

func (p *PacketFramer) Pull() (*Packet, error) {
    if p.idx < HEADER_SIZE {
        return nil, nil
    }

    if p.buf[0] != VERSION {
        return nil, PacketVersionMismatch
    }

    packetLen := getPacketLength(p.buf)
    fullLen := packetLen + HEADER_SIZE
    if packetLen == PACKET_PAYLOAD_SIZE {
        return nil, PacketMaxSizeExceeded
    }

    if fullLen <= uint16(p.idx) {
        out := make([]byte, fullLen, fullLen)
        copy(out, p.buf[:fullLen])
        copy(p.buf, p.buf[fullLen:])
        p.idx = p.idx - int(fullLen)

        pkt := PacketFromBytes(out)
        return &pkt, nil
    }

    return nil, nil
}

func FrameReader(ctx context.Context, r io.Reader, out chan *Packet) error {
    b := make([]byte, PACKET_MAX_SIZE, PACKET_MAX_SIZE)
    framer := NewPacketFramer()

    for {
        n, err := r.Read(b)
        select {
        case <-ctx.Done():
            break;
        default:
        }

        if err != nil {
            return err
        }

        fmt.Printf("LOOK: %d\n", n)
        framer.Push(b[:n])
        for {
            p, err := framer.Pull()
            if err != nil || p == nil {
                return err
            }

            out <- p
        }
    }
}
