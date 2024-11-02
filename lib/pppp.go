package pppp

import (
	"log"
	"net"
	"time"
)

// Constants
const (
	Protocol       = "udp4"
	ProbePort      = 32108
	BroadcastMsg   = "2cba5f5d"
	BroadcastDelay = 1 * time.Second
	Timeout        = 30 * time.Second
	BufferSize     = 1024
	ReadTimeout    = 5 * time.Second
)

// Connection represents P2P connection
type Connection struct {
	RemoteAddr  *net.UDPAddr
	Socket      *net.UDPConn
	isConnected chan bool
}

func InitiateConnection() {
	log.Println("Initializing Connection...")
}
