package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const num = 12   // number of nodes totally
const psnum = 3  // the shard size of process shard
const csnum = 12 // the shard size of control shard
var adv = [2]int{0, 0}
var T1 = 15 // T1:time-bound-1
// lamda
const lamda = 35 // T2 = (pontrol shard number - suceess pontrol shard)*(lamda)

var random = false // if shuffle the node by drand number
var Start time.Time

const psincs = csnum / psnum

func main() {
	// Generate nodes
	node, _ := generateNodes(num)
	getDrandResult()
	readRandomnessFromFile("./data/drand.json")
	advlist := generateadvepcoh(adv[:], num, psnum, csnum)
	for _, tempnode := range node {
		tempnode.bootstrap(random)
		for epoch, epochadvlist := range advlist {
			nodeid, _ := strconv.Atoi(tempnode.ID[4:])
			if contains(epochadvlist, nodeid) {
				tempnode.Advepoch = append(tempnode.Advepoch, epoch+1)
			}
		}
		fmt.Println(tempnode.Advepoch)
	}

	for _, tempnode := range node {
		go startNode(tempnode)
	}

	for {
		time.Sleep(2 * time.Second)
	}

}

func startNode(node *Node) {
	n := SetupNetwork(*node)
	node.PSchain = *NewPSchain(node.ID, node.PSLeader)
	node.CSchain = *NewCSchain(node.ID)
	go n.Run()
	go node.PSchain.Run(node)
	go node.CSchain.Run(&node.PSchain, node)
	time.Sleep(2 * time.Second)
	if node.PSLeader == 1 {
		go produceBlocks(node, n)
	}
	for {
		select {
		case message := <-n.IncomingMessages:
			// Processing incoming messages
			//log.Printf("%s Received message from: %+v", node.ID, message.SenderID)
			node.PSchain.AssignMessage(message, node.PSID, &node.CSchain, *node)
		case newvote := <-node.PSchain.VoteQueue_send:
			// Send vote
			newvote.Origin = node.ID
			new_vote_byte, err := EncodeVote(&newvote)
			if err != nil {
				log.Printf("error happened vote to bytes")
			}
			n.PSBroadcastMessages <- NewMessage(node.ID, 2, new_vote_byte)
			n.CSBroadcastMessages <- NewMessage(node.ID, 2, new_vote_byte)
		case blocktocs := <-node.PSchain.BlockQueue_send:
			if node.PSLeader == 1 {
				blocktocs_byte, err := EncodePSBlock(&blocktocs)
				if err != nil {
					log.Printf("error happened blocktocs to bytes")
				}
				n.CSBroadcastMessages <- NewMessage(node.ID, 1, blocktocs_byte)
			}
		case newvote := <-node.CSchain.VoteQueue_send:
			newvote.Origin = node.ID
			new_vote_byte, err := EncodeVote(&newvote)
			if err != nil {
				log.Printf("error happened vote to bytes")
			}
			n.PSBroadcastMessages <- NewMessage(node.ID, 2, new_vote_byte)
			n.CSBroadcastMessages <- NewMessage(node.ID, 2, new_vote_byte)

		}

	}
}

func produceBlocks(node *Node, n *Network) {
	sentBlockNum := 0
	// Generate the first block
	time.Sleep(40 * time.Second)
	bytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		bytes[i] = byte(rand.Intn(256))
	}
	tempBlock := NewPSBlock(1, node.ID, node.PSID, [32]byte{}, bytes, []byte("node x"))
	tempBlockByte, err := EncodePSBlock(&tempBlock)
	if err != nil {
		log.Printf("error happened block to bytes")
	}
	nodeMsg := NewMessage(node.ID, 1, tempBlockByte)
	Start = time.Now()
	n.PSBroadcastMessages <- nodeMsg
	node.PSchain.AssignMessage(nodeMsg, node.PSID, &node.CSchain, *node)
	for {
		if node.PSchain.Len > sentBlockNum {
			if node.timerExpired1 && node.timerExpired2 {
				node.timerExpired1 = false
				node.timerExpired2 = false
				// keep generate the block
				bytes := make([]byte, 256)
				for i := 0; i < 256; i++ {
					bytes[i] = byte(rand.Intn(256))
				}
				tempBlock := NewPSBlock(int16(node.PSchain.Len)+1, node.ID, node.PSID, node.PSchain.CurrentBlock.BlockHeader.Hash, bytes, []byte("node1"))
				tempBlockByte, err := EncodePSBlock(&tempBlock)
				if err != nil {
					log.Printf("error happened block to bytes")
				}
				nodeMsg := NewMessage(node.ID, 1, tempBlockByte)
				n.PSBroadcastMessages <- nodeMsg
				node.PSchain.AssignMessage(nodeMsg, node.PSID, &node.CSchain, *node)
				sentBlockNum++
			}
		}
	}
}

func (bc *PSchain) AssignMessage(msg Message, psid int, cc *CSchain, node Node) {
	switch msg.Type {
	case 1:
		temp_block, _ := DecodePSBlock(msg.Data)
		if temp_block.PSindex == psid {
			bc.BlockQueue <- *temp_block
		} else {
			cc.BlockQueue <- *temp_block
		}

	case 2:
		temp_vote, _ := DecodeVote(msg.Data)
		if temp_vote.PSindex == psid && stringInSlice(temp_vote.Origin, node.PSPeerID) {
			bc.VoteQueue <- *temp_vote
		} else {
			cc.VoteQueue <- *temp_vote
		}
	}

}
