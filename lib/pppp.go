package pppp

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Constants for network configuration
const (
	Protocol       = "udp4"
	ProbePort      = 32108
	BroadcastMsg   = "2cba5f5d"
	BroadcastDelay = 1 * time.Second
	Timeout        = 30 * time.Second
	BufferSize     = 1024
	ReadTimeout    = 5 * time.Second
)

// Connection represents a P2P connection with its associated state and channels
type Connection struct {
	RemoteAddr    *net.UDPAddr
	Socket        *net.UDPConn
	isConnected   chan bool
	PacketChannel chan Packet
	stopBroadcast chan struct{}
	mu            sync.RWMutex
	punchCount    int
	VideoHandler  *VideoHandler
}

// NewConnection creates and initializes a new Connection
func NewConnection() (*Connection, error) {
	addr := &net.UDPAddr{
		Port: 0, // Random port
		IP:   net.IPv4(0, 0, 0, 0),
	}

	conn, err := net.ListenUDP(Protocol, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP listener: %w", err)
	}

	return &Connection{
		Socket:        conn,
		isConnected:   make(chan bool),
		PacketChannel: make(chan Packet),
		stopBroadcast: make(chan struct{}),
		VideoHandler:  NewVideoHandler(),
	}, nil
}

// InitiateConnection starts the P2P connection process
func InitiateConnection() (*Connection, error) {
	log.Println("Initializing connection...")

	connection, err := NewConnection()
	if err != nil {
		return nil, err
	}

	// Start the connection handlers
	bufferChan := make(chan []byte)
	go connection.listen(bufferChan)
	go connection.broadcast()
	go connection.processPackets(bufferChan)

	// Wait for connection or timeout
	select {
	case <-connection.isConnected:
		log.Println("Connection established")
		return connection, nil
	case <-time.After(Timeout):
		connection.Close()
		return nil, fmt.Errorf("connection timed out after %v", Timeout)
	}
}

// listen continuously listens for incoming packets
func (c *Connection) listen(bufferChan chan []byte) {
	buffer := make([]byte, BufferSize)
	for {
		if err := c.Socket.SetReadDeadline(time.Now().Add(ReadTimeout)); err != nil {
			log.Printf("Failed to set read deadline: %v", err)
			continue
		}

		n, remoteAddr, err := c.Socket.ReadFromUDP(buffer)
		if err != nil {
			if !isTimeout(err) {
				log.Printf("Error reading from UDP: %v", err)
			}
			continue
		}

		c.mu.Lock()
		c.RemoteAddr = remoteAddr
		c.mu.Unlock()

		bufferCopy := make([]byte, n)
		copy(bufferCopy, buffer[:n])
		bufferChan <- bufferCopy
	}
}

// broadcast sends periodic discovery messages
func (c *Connection) broadcast() {
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: ProbePort,
	}

	broadcastMessage, err := hex.DecodeString(BroadcastMsg)
	if err != nil {
		log.Printf("Failed to decode broadcast message: %v", err)
		return
	}

	ticker := time.NewTicker(BroadcastDelay)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopBroadcast:
			return
		case <-ticker.C:
			c.mu.RLock()
			if c.punchCount > 0 {
				c.mu.RUnlock()
				return
			}
			c.mu.RUnlock()

			if _, err := c.Socket.WriteTo(broadcastMessage, broadcastAddr); err != nil {
				log.Printf("Failed to send broadcast message: %v", err)
			}
		}
	}
}

// processPackets handles incoming packets
func (c *Connection) processPackets(bufferChan chan []byte) {
	for buffer := range bufferChan {
		decryptedBuffer := decrypt(buffer)
		packet := parsePacket(decryptedBuffer)
		// log.Printf("Received packet type: %s", packet.Type)
		c.handlePacket(packet, buffer)
	}
}

// handlePacket processes different types of packets
func (c *Connection) handlePacket(packet Packet, message []byte) {
	switch packet.Type {
	case TYPE_DICT[MSG_PUNCH]:
		c.mu.Lock()
		c.punchCount++
		c.mu.Unlock()

		if err := c.Send(message); err != nil {
			log.Printf("Failed to send punch reply: %v", err)
		}
	case TYPE_DICT[MSG_P2P_RDY]:
		select {
		case c.isConnected <- true:
		default:
		}
	case TYPE_DICT[MSG_ALIVE]:
		alivePacket := prepareAlivePacket()
		if err := c.SendEncrypted(alivePacket); err != nil {
			log.Printf("Failed to send alive packet: %v", err)
		}
	case TYPE_DICT[MSG_DRW]:
		drwAckPacket := prepareDRWACKPacket(packet)
		if err := c.SendEncrypted(drwAckPacket); err != nil {
			log.Printf("Failed to send DRW ACK packet: %v", err)
		}
		if err := c.SendEncrypted(drwAckPacket); err != nil {
			log.Printf("Failed to send DRW ACK packet: %v", err)
		}
		if packet.Channel == 1 {
			c.VideoHandler.HandlePacket(packet)
		}
	case TYPE_DICT[MSG_DRW_ACK]:
		log.Println("Received: MSG_DRW_ACK")
	default:
		log.Printf("Unknown packet type received: %s", packet.Type)
	}
}

// RequestVideoStream initiates a video stream request
func (c *Connection) RequestVideoStream() chan VideoFrame {
	log.Printf("Requesting video stream from %v", c.RemoteAddr)

	streamRequest := map[string]interface{}{
		"pro":    "stream",
		"cmd":    111,
		"video":  1,
		"user":   "admin",
		"pwd":    "6666",
		"devmac": "0000",
	}

	msgBuffer, err := json.Marshal(streamRequest)

	if err != nil {
		log.Panicln("failed to marshal stream request: %w", err)
	}

	cmdPacket := prepareCommandPacket(msgBuffer)
	drwPacket := prepareDRWPacket(0, cmdPacket)

	if err := c.SendEncrypted(drwPacket); err != nil {
		log.Panicln("failed to send video stream request: %w", err)
	}

	return c.VideoHandler.VideoFrameChan
}

func (c *Connection) Send(buff []byte) error {
	_, err := c.Socket.WriteToUDP(buff, c.RemoteAddr)
	return err
}

func (c *Connection) SendEncrypted(buff []byte) error {
	return c.Send(encrypt(buff))
}

// Close cleanly shuts down the connection
func (c *Connection) Close() {
	close(c.stopBroadcast)
	if c.Socket != nil {
		c.Socket.Close()
	}
}

// isTimeout checks if an error is a timeout error
func isTimeout(err error) bool {
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout()
	}
	return false
}
