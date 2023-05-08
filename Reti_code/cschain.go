package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type CSchain struct {
	NodeID         string
	CurrentBlock   CSBlock
	Blocks         []CSBlock
	BlockQueue     chan PSBlock
	Votetable      map[int][]string
	VoteQueue      chan Vote
	VoteQueue_send chan Vote
	Len            int
}

func NewCSchain(node string) *CSchain {
	return &CSchain{
		node,
		NewCSGenesisBlock(),
		[]CSBlock{NewCSGenesisBlock()},
		make(chan PSBlock, 100),
		map[int][]string{},
		make(chan Vote, 100),
		make(chan Vote, 100),
		0,
	}
}

func SetTime2(cc *CSchain, bc *PSchain, node *Node) {
	for {
		if node.timerExpired1 && node.timerExpired2 == false {
			count := 0
			fmt.Println(cc.Votetable)
			for i := 0; i < psincs; i++ {
				if len(cc.Votetable[i]) == (psnum - 1) {
					count++
				}
			}
			if node.addps {
				count++
			}

			node.T2 = (psincs - count) * lamda
			fmt.Println(node.T2)
			Endepoch(cc, bc, node)

		}
	}

}

func Endepoch(cc *CSchain, bc *PSchain, node *Node) {
	done := make(chan bool)
	node.globalTimer2 = time.AfterFunc(time.Duration(node.T2)*time.Second, func() {
		node.timerExpired2 = true
		if node.addps == false {
			vote := cc.Votetable[node.PSID]
			if len(vote) > (csnum/2)-2 {

				bc.Blocks = append(bc.Blocks, bc.MaybeNext)
				bc.Len += 1
				fmt.Println(bc.NodeID, "--------after control shard's vote : Add Block", bc.Len)
				err := os.MkdirAll("./data/"+bc.NodeID+"/psblock", os.ModePerm)
				err = os.MkdirAll("./data/"+bc.NodeID+"/state", os.ModePerm)
				if err != nil {
				}
				tempstate := NewState(bc.CurrentBlock.Epoch, bc.CurrentBlock.Hash)
				bc.StateChain.Add(tempstate)
				SavePSBlockToFile(bc.CurrentBlock, "./data/"+bc.NodeID+"/psblock/"+strconv.Itoa(bc.Len))
				SaveStateToFile(bc.StateChain.LastState(), "./data/"+bc.NodeID+"/state/"+strconv.Itoa(bc.Len))

			}
		}

		done <- true

	})
	<-done
	node.addps = false
	cc.Add()
}

func (cc *CSchain) Add() {
	temp_block := NewCSBlock(cc.CurrentBlock.Epoch, cc.CurrentBlock.Origin, cc.CurrentBlock.PSBlock, cc.Votetable, cc.CurrentBlock.PrevHash)
	cc.Blocks = append(cc.Blocks, temp_block)
	err := os.MkdirAll("./data/"+cc.NodeID+"/csblock", os.ModePerm)
	if err != nil {
	}
	cc.Len += 1
	SaveCSBlockToFile(temp_block, "./data/"+cc.NodeID+"/csblock/"+strconv.Itoa(cc.Len))
	cc.Votetable = make(map[int][]string)
	cc.CurrentBlock = NewCSBlock(temp_block.Epoch+1, cc.NodeID, []PSBlock{}, map[int][]string{}, cc.CurrentBlock.BlockHeader.Hash)
	elapsed := time.Since(Start)
	err = saveElapsedTimeToFile(elapsed, "./data/"+cc.NodeID+"/csblock/"+strconv.Itoa(cc.Len)+"_time")
	if err != nil {
		fmt.Println("save file error:", err)
		return
	}

}

func (cc *CSchain) Run(bc *PSchain, node *Node) {
	go SetTime2(cc, bc, node)
	for {
		select {
		case temp_block := <-cc.BlockQueue:
			if temp_block.Epoch == int16(cc.Len)+1 {
				// add to current block's ps array
				cc.CurrentBlock.PSBlock = append(cc.CurrentBlock.PSBlock, temp_block)
				fmt.Println(cc.NodeID, "-----: Receive from ps", temp_block.PSindex, "block", temp_block.Epoch)
				//产生一个投票
				if !(contains(node.Advepoch, int(temp_block.Epoch))) {
					newvote := Vote{"", temp_block.Epoch, temp_block.PSindex, true, temp_block.BlockHeader.Hash, [32]byte{}}
					cc.VoteQueue_send <- newvote
				}
			}
		case temp_vote := <-cc.VoteQueue:
			if temp_vote.Epoch == int16(cc.Len)+1 {
				nowvoteRaw := cc.Votetable[temp_vote.PSindex]
				if stringInSlice(temp_vote.Origin, nowvoteRaw) {
					// Origin has already existed in this Block's Nodeid
				} else {
					nowvoteRaw = append(nowvoteRaw, temp_vote.Origin)
					cc.Votetable[temp_vote.PSindex] = nowvoteRaw
					fmt.Println(cc.NodeID, "--------Control shard :recevied ps", temp_vote.PSindex, "vote from", temp_vote.Origin, "epoch", temp_vote.Epoch)
				}
			} else {
				go cc.ccinsertvote_delay(temp_vote)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

func (cc *CSchain) ccinsertvote_delay(temp_vote Vote) {
	time.Sleep(1 * time.Second)
	cc.VoteQueue <- temp_vote
}

func (cc *CSchain) ccinsertblock_delay(temp_block PSBlock) {
	time.Sleep(1 * time.Second)
	cc.BlockQueue <- temp_block
}

func saveElapsedTimeToFile(elapsed time.Duration, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Generate Block use time: %s\n", elapsed)
	if err != nil {
		return err
	}

	return nil
}
