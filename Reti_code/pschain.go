package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type PSchain struct {
	NodeID          string
	IFPSleader      int
	CurrentBlock    PSBlock
	MaybeNext       PSBlock
	Blocks          []PSBlock
	BlockQueue      chan PSBlock
	BlockQueue_send chan PSBlock
	VoteQueue       chan Vote
	VoteQueue_send  chan Vote
	Vote4NextBlock  sync.Map
	Len             int
	StateChain      *StateChain
}

type Blockvotestatus struct {
	Num    int
	Nodeid []string
	Block  *PSBlock
}

func (bc *PSchain) AdPSBlock(b PSBlock) {
	bc.CurrentBlock = b
	bc.Blocks = append(bc.Blocks, b)
	bc.Len += 1
}

func NewPSchain(node string, IFPSleader int) *PSchain {
	return &PSchain{
		node,
		IFPSleader,
		NewPSGenesisBlock(),
		NewPSGenesisBlock(),
		[]PSBlock{NewPSGenesisBlock()},
		make(chan PSBlock, 100),
		make(chan PSBlock, 100),
		make(chan Vote, 100),
		make(chan Vote, 100),
		sync.Map{},
		0,
		NewStateChain(),
	}
}
func (bc *PSchain) Run(node *Node) {
	for {
		select {
		case temp_block := <-bc.BlockQueue:
			fmt.Println(bc.NodeID, "Now in:", (bc.Len)+1, "epoch")
			if temp_block.Epoch > int16(bc.Len) {
				if _, ok := bc.Vote4NextBlock.Load(temp_block.BlockHeader.Hash); ok {
				} else {
					var temp []string
					bc.Vote4NextBlock.Store(temp_block.BlockHeader.Hash, Blockvotestatus{1, temp, &temp_block})
					fmt.Println(bc.NodeID, "-----: Get block", temp_block.Epoch)
					node.timerExpired1 = false
					bc.MaybeNext = temp_block
					node.globalTimer1 = time.AfterFunc(time.Duration(T1)*time.Second, func() {
						node.timerExpired1 = true
						node.timerExpired2 = false
						if node.addps {
							fmt.Println(bc.NodeID, "'s epoch", temp_block.Epoch, "add process block in time-bound-1 successfully!")
						} else {
							fmt.Println(bc.NodeID, "'s epoch", temp_block.Epoch, "failed to add block in time-bound-1.")
							bc.BlockQueue_send <- temp_block
						}
					})
					fmt.Println(bc.NodeID, "Turn on time-bound-1 countdown")
					if temp_block.Origin != bc.NodeID && !contains(node.Advepoch, int(temp_block.Epoch)) {
						newvote := Vote{"", temp_block.Epoch, temp_block.PSindex, true, temp_block.BlockHeader.Hash, [32]byte{}}
						bc.VoteQueue_send <- newvote
					}
				}
			}

		case temp_vote := <-bc.VoteQueue:
			if temp_vote.Epoch == int16(bc.Len)+1 {
				if nowvoteRaw, ok := bc.Vote4NextBlock.Load(temp_vote.BlockHash); ok {
					nowvote, ok := nowvoteRaw.(Blockvotestatus)
					if !ok {
						// If the type assertion fails, the error is handled
						log.Printf("Error: Vote4NextBlock value of key %v has wrong type\n", temp_vote.BlockHash)
						continue
					}
					if stringInSlice(temp_vote.Origin, nowvote.Nodeid) {
						// Origin already exists in the Nodeid of the Block
					} else {
						nowvote.Nodeid = append(nowvote.Nodeid, temp_vote.Origin)
						bc.Vote4NextBlock.Store(temp_vote.BlockHash, Blockvotestatus{nowvote.Num + 1, nowvote.Nodeid, nowvote.Block})
						fmt.Println(bc.NodeID, "--------:recevied vote from", temp_vote.Origin)
						bc.GenerateBlock(node)
					}
				} else {
					fmt.Println(bc.NodeID, "Stuck here, please restart (server concurrency is poor）")
					go bc.insertvote_delay(temp_vote)
				}
			} else if temp_vote.Epoch > int16(bc.Len)+1 {
				fmt.Println(bc.NodeID, "receive", temp_vote.Origin, "epoch", temp_vote.Epoch, "'s vote")
				go bc.insertvote_delay(temp_vote)
			} else {

			}
		}
	}
}

func (bc *PSchain) GenerateBlock(node *Node) {
	bc.Vote4NextBlock.Range(func(key, value interface{}) bool {
		blockvotestatus, ok := value.(Blockvotestatus)
		if !ok {
			return false
		}
		if bc.IFPSleader == 1 && blockvotestatus.Num > (psnum-1) && blockvotestatus.Block.Epoch == int16(bc.Len)+1 {
			bc.Blocks = append(bc.Blocks, *blockvotestatus.Block)
			bc.CurrentBlock = *blockvotestatus.Block
			bc.Len += 1
			bc.Vote4NextBlock.Delete(key)
			node.addps = true
			fmt.Println(bc.NodeID, "--------:Add Block", bc.Len)
			err := os.MkdirAll("./data/"+bc.NodeID+"/psblock", os.ModePerm)
			err = os.MkdirAll("./data/"+bc.NodeID+"/state", os.ModePerm)
			if err != nil {
			}
			tempstate := NewState(bc.CurrentBlock.Epoch, bc.CurrentBlock.Hash)
			bc.StateChain.Add(tempstate)
			SavePSBlockToFile(bc.CurrentBlock, "./data/"+bc.NodeID+"/psblock/"+strconv.Itoa(bc.Len))
			SaveStateToFile(bc.StateChain.LastState(), "./data/"+bc.NodeID+"/state/"+strconv.Itoa(bc.Len))
		}
		if (!contains(node.Advepoch, int(blockvotestatus.Block.Epoch))) && bc.IFPSleader == 0 && blockvotestatus.Num > (psnum-2) && blockvotestatus.Block.Epoch == int16(bc.Len)+1 {
			bc.Blocks = append(bc.Blocks, *blockvotestatus.Block)
			bc.CurrentBlock = *blockvotestatus.Block
			bc.Len += 1
			bc.Vote4NextBlock.Delete(key)
			node.addps = true
			fmt.Println(bc.NodeID, "--------:Add Block", bc.Len)
			err := os.MkdirAll("./data/"+bc.NodeID+"/psblock", os.ModePerm)
			err = os.MkdirAll("./data/"+bc.NodeID+"/state", os.ModePerm)
			if err != nil {
			}
			tempstate := NewState(bc.CurrentBlock.Epoch, bc.CurrentBlock.Hash)
			bc.StateChain.Add(tempstate)
			SavePSBlockToFile(bc.CurrentBlock, "./data/"+bc.NodeID+"/psblock/"+strconv.Itoa(bc.Len))
			SaveStateToFile(bc.StateChain.LastState(), "./data/"+bc.NodeID+"/state/"+strconv.Itoa(bc.Len))
		}

		return true
	})
}

func (bc *PSchain) insertvote_delay(temp_vote Vote) {
	time.Sleep(1 * time.Second)
	bc.VoteQueue <- temp_vote
}
func (bc *PSchain) insertblock_delay(temp_block PSBlock) {
	time.Sleep(1 * time.Second)
	bc.BlockQueue <- temp_block
}
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
