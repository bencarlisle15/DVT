package main

import (
	"fmt"
	"net"
	"os"
)

func FindPort() (net.Listener, uint16){
	for port := uint16(1515); port <= uint16(65535); port += 1  {
		l, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			return l, port
		}
	}
	return nil, 0
}

func main() {
	l, port := FindPort()
	if !CheckForDatabase() {
		fmt.Println("Database could not be created")
		return
	}
	if l == nil {
		fmt.Println("Could not start server")
		return
	}
	defer l.Close()
	var self *Node = &Node{"127.0.0.1", port, 0}
	self.Id = self.CreateId()
	var node *Node = &Node{"127.0.0.1",  1515, 0}
	node.Id = node.CreateId()
	var state *State = &State{nil, nil, nil, FingerTable{make(map[int]*Node), 32, 0}}
	fmt.Printf("Starting with id %d\n", self.Id)
	state.Join(self, node)
	state.Start()
	var executor = NewExecutor(8, state)
	go executor.Start()
	fmt.Printf("Started Server on port %d\n", port)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go executor.AddConnection(c)
	}
}

func CheckForDatabase() bool {
	_, err := os.Stat("data")
	if err != nil {
		_ = os.Remove("data")
		return os.Mkdir("data", 777) == nil
	}
	return true
}