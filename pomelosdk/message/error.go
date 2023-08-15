package message

import "errors"

/**
 * ==========================
 *    Message Error Types
 * ==========================
 *
 * ErrWrongMessageType
 * ErrInvalidMessage
 * ErrRouteInfoNotFound
 *
 */
var (
	ErrWrongMessageType  = errors.New("wrong message type")
	ErrInvalidMessage    = errors.New("invalid message")
	ErrRouteInfoNotFound = errors.New("route info not found in dictionary")
)
