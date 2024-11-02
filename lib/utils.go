package pppp

import "encoding/binary"

var DRW_PACKET_INDEX uint16 = 0

func prepareCommandPacket(data []byte) []byte {
	buf := make([]byte, len(data)+8)
	buf[0] = 0x06
	buf[1] = 0x0a
	buf[2] = 0xa0
	buf[3] = 0x80
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(data)))
	copy(buf[8:], data)
	return buf
}

func prepareDRWPacket(channel int, data []byte) []byte {
	buf := make([]byte, len(data)+8)
	buf[0] = MCAM
	buf[1] = MSG_DRW
	binary.BigEndian.PutUint16(buf[2:4], uint16(len(data)+4))
	buf[4] = MDRW
	buf[5] = byte(channel)
	binary.BigEndian.PutUint16(buf[6:8], DRW_PACKET_INDEX)
	DRW_PACKET_INDEX++
	copy(buf[8:], data)
	return buf
}

func prepareAlivePacket() []byte {
	buf := make([]byte, 4)
	buf[0] = MCAM
	buf[1] = MSG_ALIVE_ACK
	binary.BigEndian.PutUint16(buf[2:4], 0)
	return buf
}

func prepareDRWACKPacket(p Packet) []byte {
	buf := make([]byte, 10)
	buf[0] = MCAM
	buf[1] = MSG_DRW_ACK
	binary.BigEndian.PutUint16(buf[2:4], 6)
	buf[4] = 0xd1
	buf[5] = p.Channel
	binary.BigEndian.PutUint16(buf[6:8], 1)
	binary.BigEndian.PutUint16(buf[8:10], p.Index)
	return buf
}
