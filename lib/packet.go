package pppp

import (
	"encoding/binary"
	"fmt"
)

const (
	MCAM          = 0xf1
	MDRW          = 0xd1
	MSG_PUNCH     = 0x41
	MSG_P2P_RDY   = 0x42
	MSG_DRW       = 0xd0
	MSG_DRW_ACK   = 0xd1
	MSG_ALIVE     = 0xe0
	MSG_ALIVE_ACK = 0xe1
	MSG_CLOSE     = 0xf0
)

var TYPE_DICT = map[byte]string{
	MSG_PUNCH:     "MSG_PUNCH",
	MSG_P2P_RDY:   "MSG_P2P_RDY",
	MSG_DRW:       "MSG_DRW",
	MSG_DRW_ACK:   "MSG_DRW_ACK",
	MSG_ALIVE:     "MSG_ALIVE",
	MSG_ALIVE_ACK: "MSG_ALIVE_ACK",
	MSG_CLOSE:     "MSG_CLOSE",
}

type Packet struct {
	Type    string
	Size    uint16
	Channel uint8
	Index   uint16
	Data    []byte
}

func (p Packet) String() string {
	return fmt.Sprintf("Packet { Type: %s, Size: %d, Channel: 0x%X, Index: %d, Data: %v }",
		p.Type, p.Size, p.Channel, p.Index, p.Data)
}

func parsePacket(buff []byte) Packet {
	var packet Packet

	if len(buff) < 8 {
		packet.Type = TYPE_DICT[buff[1]]
		return packet // Return zero-value packet if the buffer is too small
	}

	// Read the packet fields
	packet.Type = TYPE_DICT[buff[1]]
	packet.Size = binary.BigEndian.Uint16(buff[2:4])
	packet.Channel = buff[5]
	packet.Index = binary.BigEndian.Uint16(buff[6:8])
	packet.Data = buff[8:]

	return packet
}
