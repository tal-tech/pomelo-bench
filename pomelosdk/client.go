package pomelosdk

import (
	"pomelo_bench/pomelosdk/codec"
)

// Callback represents the callback type which will be called
// when the correspond events is occurred.
type Callback func(data []byte)

// NewConnector create a new Connector
func NewConnector() *Connector {
	return &Connector{
		die:    make(chan byte),
		codec:  codec.NewDecoder(),
		chSend: make(chan []byte, 64),
		//chSend:    make(chan []byte),
		mid:       1,
		events:    map[string]Callback{},
		responses: map[uint]Callback{},
	}
}
