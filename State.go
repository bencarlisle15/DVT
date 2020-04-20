package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type State struct {
	Self        *Node
	Successor   *Node
	Predecessor *Node
	FingerTable FingerTable
}

func (state *State) Create(self *Node) {
	fmt.Println("Creating ring")
	state.Self = self
	state.Predecessor = nil
	state.Successor = self
}

func (state *State) Join(self *Node, node *Node) {
	if node.Id == self.Id {
		state.Create(self)
	} else {
		response := node.FindSuccessor(self.Id)
		if response[0] == 1 {
			state.Self = self
			state.Predecessor = nil
			ip := net.IP(response[1:5]).String()
			port := binary.BigEndian.Uint16(response[5:7])
			id := binary.BigEndian.Uint64(response[7:15])
			state.Successor = &Node{ip, port, id}
			fmt.Printf("successor from %s\n", fmt.Sprintf("%s:%d", ip, port))
		} else {
			state.Create(self)
		}
	}
}

func (state *State) Stabilize() {
	if state.Predecessor != nil {
		fmt.Printf("Predecessor is %d\n", state.Predecessor.Port)
	}
	fmt.Printf("Successor is %d\n", state.Successor.Port)
	var successPredess *Node
	if state.Successor.Id == state.Self.Id && state.Predecessor != nil {
		fmt.Println("Setting succ to predecc")
		state.Successor = state.Predecessor
	} else {
		successPredess = state.Successor.GetPredecessor()
	}
	if successPredess != nil && IsBetween(state.Self.Id, successPredess.Id, state.Successor.Id, false) {
		fmt.Println("Setting succ from in between")
		state.Successor = successPredess
	}
	state.Successor.Notify(state.Self)
}

func (state *State) FixFingers() {
	state.FingerTable.NextPos %= state.FingerTable.NumFingers
	nextFinger := state.FingerTable.FindSuccessorForNextFinger(state.Self)
	state.FingerTable.Fingers[state.FingerTable.NextPos] = nextFinger
	state.FingerTable.NextPos++
}

func (state *State) CheckPredecessor() {
	if state.Predecessor != nil && state.Predecessor.IsDead() {
		state.Predecessor = nil
	}
}

func RunStabilize(state *State) {
	for {
		time.Sleep(1 * time.Second)
		state.Stabilize()
	}
}

func RunFixFingers(state *State) {
	for {
		time.Sleep(3 * time.Second)
		state.FixFingers()
	}
}

func RunCheckPredecessor(state *State) {
	for {
		time.Sleep(5 * time.Second)
		state.CheckPredecessor()
	}
}

func (state *State) Start() {
	go RunStabilize(state)
	//go RunFixFingers(state)
	//go RunCheckPredecessor(state)

}