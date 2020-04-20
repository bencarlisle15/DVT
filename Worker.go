package main

import (
	"net"
	"sync"
)

type Worker struct {
	Connections ConnectionQueue
	Lock *sync.Mutex
	Condition *sync.Cond
	State *State
}

func NewWorker(state *State) *Worker {
	var lock = &sync.Mutex{}
	worker := Worker{NewConnectionQueue(), lock, sync.NewCond(lock), state}
	return &worker
}

func (worker *Worker) Start() {
	for {
		worker.Lock.Lock()
		for worker.Connections.Len() == 0 {
			worker.Condition.Wait()
		}
		var conn = worker.Connections.Pop()
		worker.Lock.Unlock()
		worker.HandleConnection(conn)
	}
}

func (worker *Worker) AddConnection(conn net.Conn) {
	worker.Lock.Lock()
	worker.Connections.Push(conn)
	worker.Condition.Broadcast()
	worker.Lock.Unlock()
}

func (worker Worker) HandleConnection(conn net.Conn) {
	defer conn.Close()
	toWrite := RequestHandler(conn, worker.State)
	if toWrite == nil {
		return
	}
	_, _ = conn.Write(toWrite)
}

func (worker Worker) Len() int {
	return worker.Connections.Len()
}