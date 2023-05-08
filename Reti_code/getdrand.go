package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type drandResult struct {
	Randomness []byte `json:"randomness"`
	Signature  []byte `json:"signature"`
}

const drandURL = "https://api.drand.sh/public/latest"

func getDrandResult() ([]byte, error) {
	resp, err := http.Get(drandURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get drand result: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get drand result: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read drand result: %v", err)
	}

	var result drandResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal drand result: %v", err)
	}

	// Save to  drand.json
	err = os.MkdirAll("./data", os.ModePerm)
	err = ioutil.WriteFile("./data/drand.json", body, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save drand result to file: %v", err)
	}

	return result.Randomness, nil
}

func generateNodes(x int) ([]*Node, error) {
	nodes := make([]*Node, x)
	for i := 0; i < x; i++ {
		nodeID := fmt.Sprintf("node%d", i+1)
		nodeAddr := fmt.Sprintf("localhost:%d", 9000+i+1)
		nodes[i] = &Node{ID: nodeID, Addr: nodeAddr}
	}

	data := ""
	for _, node := range nodes {
		data += fmt.Sprintf("ID: %s, Addr: %s\n", node.ID, node.Addr)
	}
	err := os.MkdirAll("./data", os.ModePerm)
	err = ioutil.WriteFile("./data/nodelist_org.txt", []byte(data), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write nodes to text file: %v", err)
	}

	return nodes, nil
}
