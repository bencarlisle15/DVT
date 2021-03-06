package main

import (
	"net"
	"sync"
)

type Executor struct {
	Connections ConnectionQueue
	Workers WorkerQueue
	Lock *sync.Mutex
	Condition *sync.Cond
	State *State
}

func NewExecutor(maxThreads int, state *State) Executor {
	var workers = NewWorkerQueue()
	var worker *Worker
	for i := 0; i < maxThreads; i++ {
		worker = NewWorker(state)
		workers.Push(worker)
		go worker.Start()
	}
	var mutex = sync.Mutex{}
	executor := Executor{NewConnectionQueue(), workers, &mutex, sync.NewCond(&mutex), state}
	return executor
}

func (executor *Executor) Start() {
	for {
		executor.Lock.Lock()
		for executor.Connections.Len() == 0 {
			executor.Condition.Wait()
		}
		var conn = executor.Connections.Pop()
		executor.Lock.Unlock()
		//fmt.Println("Incoming Connection")
		var smallestWorker = SmallestWorker(executor.Workers)
		go smallestWorker.AddConnection(conn)
	}
}

func (executor *Executor) AddConnection(conn net.Conn) {
	executor.Lock.Lock()
	executor.Connections.Push(conn)
	executor.Condition.Broadcast()
	executor.Lock.Unlock()
}

func SmallestWorker(queue WorkerQueue) *Worker {
	var smallestWorker = queue.Get(0)
	for i := 1; i < queue.Len(); i++ {
		if smallestWorker.Len() > queue.Get(i).Len() {
			smallestWorker = queue.Get(i)
		}
	}
	return smallestWorker
}