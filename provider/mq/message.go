package mq

import (
	"context"
	"sync"
)

// Payload is the Message's payload.
type Payload []byte

type ackType int

const (
	noAckSent ackType = iota
	ack
	nack
)

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

// Message is the basic transfer unit.
// Messages are emitted by Publishers and received by Subscribers.
type Message struct {
	// Payload is the message's payload.
	Payload Payload

	// ack is closed, when acknowledge is received.
	ack chan struct{}
	// noACk is closed, when negative acknowledge is received.
	noAck chan struct{}

	ackMutex    sync.Mutex
	ackSentType ackType

	ctx context.Context
}

// NewMessage creates a new Message with payload.
func NewMessage(payload Payload) *Message {
	return &Message{
		Payload: payload,
		ack:     make(chan struct{}),
		noAck:   make(chan struct{}),
	}
}

// Ack sends message's acknowledgement.
//
// Ack is not blocking.
// Ack is idempotent.
// False is returned, if Nack is already sent.
func (m *Message) Ack() bool {
	m.ackMutex.Lock()
	defer m.ackMutex.Unlock()

	if m.ackSentType == nack {
		return false
	}
	if m.ackSentType != noAckSent {
		return true
	}

	m.ackSentType = ack
	if m.ack == nil {
		m.ack = closedchan
	} else {
		close(m.ack)
	}

	return true
}

// Nack sends message's negative acknowledgement.
//
// Nack is not blocking.
// Nack is idempotent.
// False is returned, if Ack is already sent.
func (m *Message) Nack() bool {
	m.ackMutex.Lock()
	defer m.ackMutex.Unlock()

	if m.ackSentType == ack {
		return false
	}
	if m.ackSentType != noAckSent {
		return true
	}

	m.ackSentType = nack

	if m.noAck == nil {
		m.noAck = closedchan
	} else {
		close(m.noAck)
	}

	return true
}

// Acked returns channel which is closed when acknowledgement is sent.
//
// Usage:
//
//	select {
//	case <-message.Acked():
//		// ack received
//	case <-message.Nacked():
//		// nack received
//	}
func (m *Message) Acked() <-chan struct{} {
	return m.ack
}

// Nacked returns channel which is closed when negative acknowledgement is sent.
//
// Usage:
//
//	select {
//	case <-message.Acked():
//		// ack received
//	case <-message.Nacked():
//		// nack received
//	}
func (m *Message) Nacked() <-chan struct{} {
	return m.noAck
}

// Context returns the message's context. To change the context, use
// SetContext.
//
// The returned context is always non-nil; it defaults to the
// background context.
func (m *Message) Context() context.Context {
	if m.ctx != nil {
		return m.ctx
	}
	return context.Background()
}

// SetContext sets provided context to the message.
func (m *Message) SetContext(ctx context.Context) {
	m.ctx = ctx
}
