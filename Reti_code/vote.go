package main

import (
	"bytes"
	"encoding/gob"
)

type Vote struct {
	Origin    string
	Epoch     int16
	PSindex   int
	Agree     bool
	BlockHash [32]byte
	Signature [32]byte
}

func NewVote(Org string, epoch int16, PSindex int, agree bool, BlockSign [32]byte) Vote {
	var sign [32]byte
	return Vote{Org, epoch, PSindex, agree, BlockSign, sign}
}

func EncodeVote(v *Vote) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeVote(data []byte) (*Vote, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var v Vote
	err := dec.Decode(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
