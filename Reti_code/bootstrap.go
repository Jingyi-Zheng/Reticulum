package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strings"
)

func (node *Node) bootstrap(random bool) {
	nodes_org := readNodesFromTextFile("./data/nodelist_org.txt")
	nodes := shuffle(nodes_org, random)
	mypsindex, mycsindex, ifpsleader, ifcsleader := getShardIndex(nodes, *node, psnum, csnum)
	node.PSID = mypsindex
	node.CSID = mycsindex
	node.CSLeader = ifcsleader
	node.PSLeader = ifpsleader
	nodeinsameps := getNodeInPS(nodes, mypsindex, psnum)
	for _, i := range nodeinsameps {
		if i.ID != node.ID {
			node.PSPeerID = append(node.PSPeerID, i.ID)
			node.PSPeerAddr = append(node.PSPeerAddr, i.Addr)
		}
	}
	nodeinsamecs := getNodeInCS(nodes, mycsindex, csnum)
	for _, i := range nodeinsamecs {
		if i.ID != node.ID && !stringInSlice(i.ID, node.PSPeerID) {
			node.CSPeerID = append(node.CSPeerID, i.ID)
			node.CSPeerAddr = append(node.CSPeerAddr, i.Addr)
		}
	}

}

func shuffle(nodes []Node, random bool) []Node {
	if random == true {
		randomness, err := readRandomnessFromFile("./data/drand.json")
		if err != nil {
			panic(err)
		}
		rand.Seed(int64(randomness[0]))
		for i := len(nodes) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
	}
	return nodes
}

func getShardIndex(nodes []Node, node Node, psSize, csSize int) (int, int, int, int) {
	nodeIndex := -1
	csleader := 0
	psleader := 0
	for i := 0; i < len(nodes); i++ {
		if nodes[i].ID == node.ID {
			nodeIndex = i
			break
		}
	}
	if nodeIndex == -1 {
		panic(fmt.Sprintf("node %s not found in the node list", node.ID))
	}
	csIndex := int(math.Floor(float64(nodeIndex) / float64(csSize)))
	if csIndex*csSize == nodeIndex {
		csleader = 1
	}
	psIndex := int(math.Floor(float64(nodeIndex) / float64(psSize)))
	if psIndex*psSize == nodeIndex {
		psleader = 1
	}
	return psIndex, csIndex, psleader, csleader
}

func getNodeInPS(nodes []Node, psIndex, psSize int) []Node {
	startIndex := psIndex * psSize
	endIndex := (psIndex + 1) * psSize
	if endIndex >= len(nodes) {
		endIndex = len(nodes)
	}
	return nodes[startIndex:endIndex]
}

func getNodeInCS(nodes []Node, csIndex, csSize int) []Node {
	startIndex := csIndex * csSize
	endIndex := (csIndex + 1) * csSize
	if endIndex >= len(nodes) {
		endIndex = len(nodes)
	}
	return nodes[startIndex:endIndex]
}

func readNodesFromTextFile(filename string) []Node {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	nodes := make([]Node, 0)

	for _, line := range lines {
		if line == "" { // check empty line
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 4 {
			nodeID := strings.TrimSuffix(fields[1], ",")
			nodeAddr := strings.TrimSuffix(fields[3], ",")
			nodes = append(nodes, Node{ID: nodeID, Addr: nodeAddr})
		}
	}

	return nodes
}

func readRandomnessFromFile(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	var result struct {
		Randomness string `json:"randomness"`
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return result.Randomness, nil
}
