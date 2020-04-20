package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"net"
)

type Node struct {
	Ip string
	Port uint16
	Id uint64
}

func (node *Node) CreateId() uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%d", node.Ip, node.Port)))
	idBytes := hasher.Sum(nil)[:8]
	return binary.BigEndian.Uint64(idBytes)
}

func (node *Node) FindSuccessor(id uint64) []byte {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", node.Ip, node.Port))
	response := make([]byte, 15)
	if err != nil {
		return response
	}
	defer conn.Close()
	request := make([]byte, 9)
	binary.BigEndian.PutUint64(request[1:9], id)
	_, _ = conn.Write(request)
	_, err = conn.Read(response)
	if err != nil {
		return response
	}
	return response
}

func (node *Node) GetPredecessor() *Node {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", node.Ip, node.Port))
	response := make([]byte, 15)
	if err != nil {
		return nil
	}
	defer conn.Close()
	request := []byte{2}
	_, _ = conn.Write(request)
	_, err = conn.Read(response)
	if err != nil || response[0] == 0 {
		return nil
	}
	ip := net.IP(response[1:5]).String()
	port := binary.BigEndian.Uint16(response[5:7])
	id := binary.BigEndian.Uint64(response[7:15])
	return &Node{ip, port, id}
}

func (node *Node) Notify(self *Node) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", node.Ip, node.Port))
	if err != nil {
		return
	}
	defer conn.Close()
	request := make([]byte, 15)
	request[0] = 1
	copy(request[1:5], net.ParseIP(self.Ip).To4())
	binary.BigEndian.PutUint16(request[5:7], self.Port)
	binary.BigEndian.PutUint64(request[7:15], self.Id)
	_, _ = conn.Write(request)
}

func (node *Node) IsDead() bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", node.Ip, node.Port))
	response := make([]byte, 1)
	if err != nil {
		return true
	}
	defer conn.Close()
	request := []byte{3}
	_, _ = conn.Write(request)
	_, err = conn.Read(response)
	return err != nil || response[0] == 0
}