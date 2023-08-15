package packet

import (
	"fmt"
)

// [Reference](https://github.com/NetEase/pomelo/wiki/Communication-Protocol)

// New --
func New() *Packet {
	return &Packet{}
}

// Packet --
type Packet struct {
	Type   byte //
	Length int  // body Content length, a big-endian integer of 3 bytes, so the maximum packet length is 2^24 bytes。
	Data   []byte
}

// String --
func (p *Packet) String() string {
	return fmt.Sprintf("Type: %d, Length: %d, Data: %s", p.Type, p.Length, string(p.Data))
}
