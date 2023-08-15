package packet

// [Reference](https://github.com/NetEase/amoeba/wiki/Communication-Protocol)

/**
 * ==========================
 *        Packet
 * ==========================
 *
 * Handshake    01 - Client-to-server handshake request and server-to-client handshake response
 * HandshakeAck 02 - Client to server handshake ack
 * Heartbeat    03 - Heartbeat packet
 * Data         04 - Data packet
 * Kick         05 - Server active disconnect notification
 *
 */
const (
	Handshake    = 0x01
	HandshakeAck = 0x02
	Heartbeat    = 0x03
	Data         = 0x04
	Kick         = 0x05
)
