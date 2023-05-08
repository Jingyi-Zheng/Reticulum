package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type State struct {
	StateHeader
	Data []byte
}

type StateChain []State

type StateHeader struct {
	Epoch        int16
	PreBlockHash [32]byte
	PrevHash     [32]byte
	Hash         [32]byte
}

func NewStateChain() *StateChain {
	chain := &StateChain{}
	chain.Add(NewState(0, [32]byte{}))
	return chain
}

func NewState(epoch int16, preblockhash [32]byte) State {
	data := []byte{}
	header := StateHeader{Epoch: epoch, PreBlockHash: preblockhash}
	return State{header, data}
}

func EncodeStateheader(sh *StateHeader) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(sh); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (chain *StateChain) Add(state State) error {
	//calculate the new hash
	lastState := chain.LastState()
	state.StateHeader.PrevHash = lastState.StateHeader.Hash
	statebytes, err := EncodeState(&state)
	if err != nil {
		return err
	}
	state.StateHeader.Hash = sha256.Sum256(statebytes)

	// add new state to chain
	*chain = append(*chain, state)
	return nil
}

// get the latest state
func (chain StateChain) LastState() State {
	if len(chain) == 0 {
		return NewState(0, [32]byte{})
	}
	return (chain)[len(chain)-1]
}

// Trun state structure into bytes
func EncodeState(state *State) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(state); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SaveStateToFile(state State, filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
	}

	// turn state into json
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	// write json to file
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	fmt.Println("State saved to file ", filename)
	return nil
}
