package pppp

import (
	"bytes"
	"encoding/hex"
	"sort"
)

type VideoHandler struct {
	videoOverflow   bool
	lastVideoFrame  int
	videoBoundaries map[uint16]bool
	videoReceived   map[uint16][]byte
	VideoFrameChan  chan VideoFrame
}

type VideoFrame struct {
	Frame       []byte
	PacketIndex uint16
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{
		videoOverflow:   false,
		lastVideoFrame:  -1,
		videoBoundaries: make(map[uint16]bool),
		videoReceived:   make(map[uint16][]byte),
		VideoFrameChan:  make(chan VideoFrame, 100),
	}
}

func (v *VideoHandler) HandlePacket(p Packet) {
	if p.Channel != 1 {
		return
	}

	v.checkAndHandleOverflow(p.Index)
	v.storePacketData(p)

	if len(v.videoBoundaries) > 1 {
		v.extractCompleteVideoFrame()
	}
}

// checkAndHandleOverflow manages packet index overflow and resets frame data as needed.
func (v *VideoHandler) checkAndHandleOverflow(index uint16) {
	if index > 65400 {
		v.videoOverflow = true
	}

	if v.videoOverflow && index < 65400 {
		v.resetFrameData()
	}
}

// storePacketData stores packet data and marks boundaries if the header matches.
func (v *VideoHandler) storePacketData(p Packet) {
	const header = "55aa15a80300"
	headerBytes, _ := hex.DecodeString(header)

	if bytes.HasPrefix(p.Data, headerBytes) {
		v.videoReceived[p.Index] = p.Data[0x20:]
		v.videoBoundaries[p.Index] = true
	} else {
		v.videoReceived[p.Index] = p.Data
	}
}

// extractCompleteVideoFrame checks for complete video frames, sends them, and cleans up old data.
func (v *VideoHandler) extractCompleteVideoFrame() {
	indices := v.sortedBoundaryIndices()
	index := indices[len(indices)-2]
	lastIndex := indices[len(indices)-1]

	if uint16(v.lastVideoFrame) == index {
		return
	}

	frameData, isComplete := v.buildFrameData(index, lastIndex)
	if isComplete {
		v.sendFrame(index, frameData)
		v.videoBoundaries = make(map[uint16]bool)
	}
}

// sortedBoundaryIndices returns a sorted list of boundary indices.
func (v *VideoHandler) sortedBoundaryIndices() []uint16 {
	var indices []uint16
	for idx := range v.videoBoundaries {
		indices = append(indices, idx)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })
	return indices
}

// buildFrameData constructs the video frame data from received packets within the specified range.
func (v *VideoHandler) buildFrameData(start, end uint16) ([]byte, bool) {
	var frameData [][]byte
	var completeness string
	var complete bool = true
	for i := start; i < end; i++ {
		if data, exists := v.videoReceived[i]; exists {
			completeness += "x"
			frameData = append(frameData, data)
		} else {
			complete = false
			completeness += "-"
		}
	}

	mergedFrame := bytes.Join(frameData, nil)

	return mergedFrame, complete
}

// sendFrame pushes the complete video frame to the VideoFrameChan.
func (v *VideoHandler) sendFrame(index uint16, frame []byte) {
	v.lastVideoFrame = int(index)
	v.VideoFrameChan <- VideoFrame{
		Frame:       frame,
		PacketIndex: index,
	}
}

// resetFrameData resets the state of video frame data when overflow conditions are met.
func (v *VideoHandler) resetFrameData() {
	v.lastVideoFrame = -1
	v.videoOverflow = false
	v.videoBoundaries = make(map[uint16]bool)
	v.videoReceived = make(map[uint16][]byte)
}
