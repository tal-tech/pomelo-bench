package pomelosdk

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"pomelo_bench/pomelosdk/codec"
	"pomelo_bench/pomelosdk/message"
	"pomelo_bench/pomelosdk/packet"
)

type (
	// Connector is a Amoeba [nano] client
	Connector struct {
		conn              *websocket.Conn // low-level connection
		codec             *codec.Decoder  // decoder
		mid               uint            // message id
		muConn            sync.RWMutex
		connecting        bool        // connection status
		die               chan byte   // connector close channel
		chSend            chan []byte // send queue
		connectedCallback func()

		// some packet data
		handshakeData    []byte // handshake data
		handshakeAckData []byte // handshake ack data
		heartbeatData    []byte // heartbeat data

		// events handler
		sync.RWMutex
		events map[string]Callback

		// response handler
		muResponses sync.RWMutex
		responses   map[uint]Callback

		heartbeat int //  = 13
	}
	// DefaultACK --
	DefaultHandshakePacket struct {
		Code int              `json:"code"`
		Sys  HeartbeatSysOpts `json:"sys"`
		Ack  int              `json:"ack"`
	}
	// HeartbeatSysOpts --
	HeartbeatSysOpts struct {
		Heartbeat int `json:"heartbeat"`
	}

	// SysOpts --
	SysOpts struct {
		Version string                 `json:"version"`
		Type    string                 `json:"type"`
		RSA     map[string]interface{} `json:"rsa"`
	}

	// HandshakeOpts --
	HandshakeOpts struct {
		Sys      SysOpts                `json:"sys"`
		UserData map[string]interface{} `json:"user"`
	}
)

// SetHandshake --
func (c *Connector) SetHandshake(handshake interface{}) error {
	data, err := json.Marshal(handshake)
	if err != nil {
		return err
	}

	c.handshakeData, err = codec.Encode(packet.Handshake, data)
	if err != nil {
		return err
	}

	return nil
}

// SetHandshakeAck --
func (c *Connector) SetHandshakeAck(handshakeAck interface{}) error {
	var err error
	if handshakeAck == nil {
		c.handshakeAckData, err = codec.Encode(packet.HandshakeAck, nil)
		if err != nil {
			return err
		}
		return nil
	}

	data, err := json.Marshal(handshakeAck)
	if err != nil {
		return err
	}

	c.handshakeAckData, err = codec.Encode(packet.HandshakeAck, data)
	if err != nil {
		return err
	}

	return nil
}

// SetHeartBeat --
func (c *Connector) SetHeartBeat(heartbeat interface{}) error {
	var err error
	if heartbeat == nil {
		c.heartbeatData, err = codec.Encode(packet.Heartbeat, nil)
		if err != nil {
			return err
		}
		return nil
	}
	data, err := json.Marshal(heartbeat)
	if err != nil {
		return err
	}

	c.heartbeatData, err = codec.Encode(packet.Heartbeat, data)
	if err != nil {
		return err
	}

	return nil
}

// Connected --
func (c *Connector) Connected(cb func()) {
	//fmt.Println("cb func is ", cb)
	c.connectedCallback = cb
}

// InitReqHandshake --
// func (c *Connector) InitReqHandshake(opts *HandshakeOpts) error {
// 	return c.SetHandshake(opts)
// }

// InitReqHandshake --
func (c *Connector) InitReqHandshake(version, hType string, rsa, userData map[string]interface{}) error {
	return c.SetHandshake(&HandshakeOpts{
		Sys: SysOpts{
			Version: version,
			Type:    hType,
			RSA:     rsa,
		},
		UserData: userData,
	})
}

// InitHandshakeACK --
func (c *Connector) InitHandshakeACK(heartbeatDuration int) error {
	ackDataMap := &DefaultHandshakePacket{
		Code: 200,
		Sys: HeartbeatSysOpts{
			Heartbeat: heartbeatDuration,
		},
		Ack: 1,
	}
	return c.SetHandshakeAck(ackDataMap)
}

// Run --
func (c *Connector) Run(ctx context.Context, addr string, tickrate int64) error {
	if c.handshakeData == nil {
		return errors.New("handshake not defined")
	}

	if c.handshakeAckData == nil {
		err := c.SetHandshakeAck(nil)
		if err != nil {
			return err
		}
	}

	if c.heartbeatData == nil {
		err := c.SetHeartBeat(nil)
		if err != nil {
			return err
		}
	}

	if c.heartbeat == 0 {
		c.heartbeat = 13
	}

	if strings.HasPrefix(addr, "ws://") || strings.HasPrefix(addr, "wss://") {

		//conn, err = websocket.Dial(addr, "", "http://localhost/")
		socketConn, _, err := websocket.DefaultDialer.DialContext(ctx, addr, nil)
		if err != nil {
			return err
		}

		c.conn = socketConn

	} else {
		return errors.New("invalid websocket")
	}

	// var err error
	// var conn net.Conn
	//conn, err := websocket.Dial(addr)

	c.connecting = true

	go c.write()

	c.send(c.handshakeData)

	err := c.read(tickrate)

	return err
}

// Request send a request to server and register a callbck for the response
func (c *Connector) Request(route string, data []byte, callback Callback) error {
	msg := &message.Message{
		Type:  message.Request,
		Route: route,
		ID:    c.mid,
		Data:  data,
	}
	//fmt.Println("request", msg)
	c.setResponseHandler(c.mid, callback)
	if err := c.sendMessage(msg); err != nil {
		c.setResponseHandler(c.mid, nil)
		return err
	}

	return nil
}

// Notify send a notification to server
func (c *Connector) Notify(route string, data []byte) error {
	msg := &message.Message{
		Type:  message.Notify,
		Route: route,
		Data:  data,
	}
	return c.sendMessage(msg)
}

// On add the callback for the event
func (c *Connector) On(event string, callback Callback) {
	c.Lock()
	defer c.Unlock()

	c.events[event] = callback
}

// Close the connection, and shutdown the benchmark
func (c *Connector) Close() {
	if !c.connecting {
		return
	}
	c.connecting = false
	c.die <- 1
	c.conn.Close()
}

// IsClosed check the connection is closed
func (c *Connector) IsClosed() bool {
	return !c.connecting
}

func (c *Connector) eventHandler(event string) (Callback, bool) {
	c.RLock()
	defer c.RUnlock()

	cb, ok := c.events[event]
	return cb, ok
}

func (c *Connector) responseHandler(mid uint) (Callback, bool) {
	c.muResponses.RLock()
	defer c.muResponses.RUnlock()

	cb, ok := c.responses[mid]
	return cb, ok
}

func (c *Connector) setResponseHandler(mid uint, cb Callback) {
	c.muResponses.Lock()
	defer c.muResponses.Unlock()

	if cb == nil {
		delete(c.responses, mid)
	} else {
		c.responses[mid] = cb
	}
}

func (c *Connector) sendMessage(msg *message.Message) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}
	//log.Printf("%+v | %+v | %+v\n", msg.Data, msg, data)

	payload, err := codec.Encode(packet.Data, data)
	if err != nil {
		return err
	}

	//log.Printf("payload %d | %+v \n", len(payload), payload)

	c.mid++

	err = c.sendWithTimeout(payload, time.Second)
	return err
}

func (c *Connector) write() {

	var heartbeat = time.After(time.Duration(c.heartbeat) * time.Second)

	for {
		select {
		case data := <-c.chSend:
			if c.conn != nil {
				if err := c.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					//log.Printf("conn write len %d , err %s ", n, err.Error())
					// c.Close()
				} else {
					//log.Println("conn write success,n: ", n, ", data:", string(data))
				}
			}

		case <-heartbeat:
			_ = c.conn.WriteMessage(websocket.BinaryMessage, c.heartbeatData)

			heartbeat = time.After(time.Duration(c.heartbeat) * time.Second)

			//log.Println("send heartbeat")

		case <-c.die:
			return
		}
	}
}

func (c *Connector) send(data []byte) {
	c.chSend <- data
}

func (c *Connector) sendWithTimeout(data []byte, timeout time.Duration) error {
	select {
	case c.chSend <- data:
		return nil
	case <-time.After(timeout):
		return errors.New("time out")
	}
}

func (c *Connector) read(tickrate int64) error {

	for {
		time.Sleep(time.Second / time.Duration(tickrate))
		if c.IsClosed() {
			return errors.New("read err: connector is closed")
		}
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("connector read err", err.Error())
			c.Close()
			return err
			// continue
		}

		packets, err := c.codec.Decode(p)
		if err != nil {
			log.Println("connector read decode err", err.Error())
			// c.Close()
			// return
			continue
		}

		for i := range packets {
			p := packets[i]
			// log.Println("packet-->", p)
			c.processPacket(p)
		}
	}
}

func (c *Connector) processPacket(p *packet.Packet) {
	// log.Printf("packet: %+v\n", p)
	switch p.Type {
	case packet.Handshake:
		var handshakeResp DefaultHandshakePacket
		err := json.Unmarshal(p.Data, &handshakeResp)
		if err != nil {
			c.Close()
			return
		}
		if handshakeResp.Code == 200 {

			if handshakeResp.Sys.Heartbeat != 0 {
				c.heartbeat = handshakeResp.Sys.Heartbeat
			}

			c.send(c.handshakeAckData)
			if c.connectedCallback != nil {
				c.connectedCallback()
			}
		} else {
			log.Println("bad packet handshake code, not 200:", string(p.Data))
			c.Close()
		}
	case packet.Data:
		msg, err := message.Decode(p.Data)
		if err != nil {
			return
		}
		c.processMessage(msg)

	case packet.Kick:
		log.Println("server kick -->", p)
		c.Close()
	}
}

func (c *Connector) processMessage(msg *message.Message) {
	switch msg.Type {
	case message.Push:
		cb, ok := c.eventHandler(msg.Route)
		if !ok {
			log.Println("event handler not found", msg.Route, msg, c.events)
			return
		}
		cb(msg.Data)

	case message.Response:
		cb, ok := c.responseHandler(msg.ID)
		if !ok {
			log.Println("response handler not found", msg.ID)
			return
		}

		cb(msg.Data)
		c.setResponseHandler(msg.ID, nil)
	}
}
