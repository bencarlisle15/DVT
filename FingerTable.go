package main

import (
	"encoding/binary"
	"net"
)

type FingerTable struct {
	Fingers map[int]*Node
	NumFingers int
	NextPos int
}

func (fingerTable *FingerTable) GetClosestPrecedingNode(self *Node, id uint64) *Node {
	for i := fingerTable.NumFingers - 1; i >= 0; i-- {
		if fingerTable.Fingers[i] != nil && IsBetween(self.Id, fingerTable.Fingers[i].Id, id, false) {
			return fingerTable.Fingers[i]
		}
	}
	return self
}

func (fingerTable *FingerTable) FindSuccessorForNextFinger(self *Node) *Node {
	id := AddPow(self.Id, fingerTable.NextPos)
	node := fingerTable.GetClosestPrecedingNode(self, id)
	if node.Id == self.Id {
		return self
	}
	response := node.FindSuccessor(id)
	if response[0] == 0 {
		return nil
	}
	ip := net.IP(response[1:5]).String()
	port := binary.BigEndian.Uint16(response[5:7])
	successorId := binary.BigEndian.Uint64(response[7:15])
	return &Node{ip, port, successorId}
}