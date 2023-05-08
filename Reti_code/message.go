package main

import (
	"bytes"
	"encoding/gob"
)

type Message struct {
	SenderID string
	Type     int16 //1:block 2:Vote
	Data     []byte
}

func NewMessage(SenderID string, Type int16, Data []byte) Message {
	return Message{SenderID, Type, Data}
}

// turn Message structure into bytes
func EncodeMsg(m *Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// turn bytes into Message
func DecodeMsg(data []byte) (*Message, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var m Message
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
