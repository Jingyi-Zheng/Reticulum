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

	"gopkg.in/yaml.v2"
)

type PSBlock struct {
	BlockHeader
	PSindex int
	Data    []byte
}
type CSBlock struct {
	BlockHeader
	Votetable map[int][]string `yaml:"votetable" json:"votetable"`
	PSBlock   []PSBlock
}

type BlockHeader struct {
	Epoch    int16
	Origin   string
	PrevHash [32]byte
	Hash     [32]byte
}

func NewPSGenesisBlock() PSBlock {
	return PSBlock{
		BlockHeader{
			Epoch:    0,
			Origin:   "",
			PrevHash: [32]byte{},
			Hash:     [32]byte{},
		},
		0,
		[]byte{},
	}
}

func NewCSGenesisBlock() CSBlock {
	return CSBlock{
		BlockHeader{
			Epoch:    0,
			Origin:   "",
			PrevHash: [32]byte{},
			Hash:     [32]byte{},
		},
		map[int][]string{},
		[]PSBlock{},
	}
}

func NewPSBlock(epoch int16, Org string, PSindex int, previousBlock [32]byte, data []byte, sign []byte) PSBlock {
	header := BlockHeader{Epoch: epoch, Origin: Org, PrevHash: previousBlock}
	block := PSBlock{
		BlockHeader: header,
		PSindex:     PSindex,
		Data:        data,
	}
	blockbytes, _ := EncodePSBlock(&block)
	block.BlockHeader.Hash = sha256.Sum256(blockbytes)
	return block
}

func NewCSBlock(epoch int16, origin string, psblock []PSBlock, vote map[int][]string, prevHash [32]byte) CSBlock {
	blockHeader := BlockHeader{
		Epoch:    epoch,
		Origin:   origin,
		PrevHash: prevHash,
	}
	block := CSBlock{
		BlockHeader: blockHeader,
		PSBlock:     psblock,
		Votetable:   vote,
	}
	blockbytes, _ := EncodeCSBlock(&block)
	block.BlockHeader.Hash = sha256.Sum256(blockbytes)
	return block
}

func EncodeCSBlock(block *CSBlock) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(block)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// turn Block structure into bytes
func EncodePSBlock(b *PSBlock) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// turn bytes to Block structure
func DecodePSBlock(data []byte) (*PSBlock, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var b PSBlock
	err := dec.Decode(&b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// turn BlockHeader structure into bytes
func encodeBlockHeader(h *BlockHeader) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(h)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 将字节码解码为 BlockHeader 结构体
func decodeBlockHeader(data []byte) (*BlockHeader, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var h BlockHeader
	err := dec.Decode(&h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
func SavePSBlockToFile(block PSBlock, filename string) error {
	// check if the file is exist or create an empty one
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
	}

	// turn byte into json
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	// write json to file
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	fmt.Println("process Block saved to file ", filename)
	return nil
}

// SaveCSBlockToFile saves the CSBlock struct to a YAML file.
func SaveCSBlockToFile(block CSBlock, filename string) error {
	// Check if the file exists, if not create an empty file
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
	}

	// Encode the block to YAML
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	data, err := yaml.Marshal(block)
	if err != nil {
		return err
	}

	// Write the YAML to the file
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	fmt.Println("Control Block saved to file", filename)
	return nil
}
