package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func RequestHandler(conn net.Conn, state *State) []byte {
	requestType := make([]byte, 1)
	_, err := conn.Read(requestType)
	if err != nil {
		return nil
	}
	//fmt.Printf("Recieved message %d\n", requestType)
	var response []byte
	switch requestType[0] {
	case 0:
		response = FindSuccessor(conn, state)
	case 1:
		response = Notify(conn, state)
	case 2:
		response = GetPredecessor(state)
	case 3:
		response = Alive()
	}
	return response
}

func FindSuccessor(conn net.Conn, state *State) []byte {
	request := make([]byte, 8)
	_, err := conn.Read(request)
	if err != nil {
		return nil
	}
	id := binary.BigEndian.Uint64(request[0:8])
	response := make([]byte, 15)
	if IsBetween(state.Self.Id, id, state.Successor.Id, true) {
		response[0] = 1
		copy(response[1: 5], net.ParseIP(state.Successor.Ip).To4())
		binary.BigEndian.PutUint16(response[5:7], state.Successor.Port)
		binary.BigEndian.PutUint64(response[7:15], state.Successor.Id)
		return response
	}
	closest := state.FingerTable.GetClosestPrecedingNode(state.Self, id)
	if closest.Id != state.Self.Id {
		response = closest.FindSuccessor(id)
	} else {
		response[0] = 1
		copy(response[1: 5], net.ParseIP(state.Successor.Ip).To4())
		binary.BigEndian.PutUint16(response[5:7], state.Successor.Port)
		binary.BigEndian.PutUint64(response[7:15], state.Successor.Id)
		return response
	}
	return response
}

func Notify(conn net.Conn, state *State) []byte {
	request := make([]byte, 14)
	_, err := conn.Read(request)
	if err != nil {
		return nil
	}
	id := binary.BigEndian.Uint64(request[6:14])
	if state.Predecessor == nil || state.Predecessor.Id == state.Self.Id || IsBetween(state.Predecessor.Id, id, state.Self.Id, false) {
		ip := net.IP(request[0:4]).String()
		port := binary.BigEndian.Uint16(request[4:6])
		state.Predecessor = &Node{ip, port, id}
		fmt.Printf("Am notified of new predess %d\n", port)
	}
	return nil
}

func GetPredecessor(state *State) []byte {
	response := make([]byte, 15)
	if state.Predecessor != nil {
		response[0] = 1
		copy(response[1:5], net.ParseIP(state.Successor.Ip).To4())
		binary.BigEndian.PutUint16(response[5:7], state.Successor.Port)
		binary.BigEndian.PutUint64(response[7:15], state.Successor.Id)
	}
	return response
}

func Alive() []byte {
	return []byte{1}
}