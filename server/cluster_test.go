package main

import (
	"testing"
	"sync/atomic"
	"github.com/amolb89/cluster/clust"
	"time"
)

func TestClust(t *testing.T) {
	configPath := "/home/amol/Desktop/zmqp/src/github.com/amolb89/cluster/server/serverConfig.json"
	server := new(clust.Serv)
	server.Set(5556,configPath)
	var sendcount uint64 = 0
	go func() {
	     for j := 1; j<1000 ;j++ {	
		//Send point to point messages
		for i:=5557 ; i<=5559; i++ {
			server.Outbox() <- &clust.Envelope{i,"hello"}
			sendcount++
		}
		//send broadcase messages
		server.Outbox() <- &clust.Envelope{-1,"hello"}
		sendcount = sendcount + 3
	     }
	}()
	
	//Server 1
	server1 := new(clust.Serv)
	server1.Set(5559,configPath)

	//Server 2
	server2 := new(clust.Serv)
	server2.Set(5557,configPath)

	//Server 3
	server3 := new(clust.Serv)
	server3.Set(5558,configPath)
	
	var ReceiveCount uint64 = 0
	go func() {
		for {
			<- server1.Inbox()
			atomic.AddUint64(&ReceiveCount, 1)
		}
	}()

	go func() {
		for {
			<- server2.Inbox()
			atomic.AddUint64(&ReceiveCount, 1)
		}
	}()

	go func() {
		for {
			<- server3.Inbox()
			atomic.AddUint64(&ReceiveCount, 1)
		}
	}()
	time.Sleep(time.Second*10)
	if sendcount != ReceiveCount {
		t.Errorf("Sending and receiving count do not match")
	}
}

func TestBroadcast(t *testing.T) {
	configPath := "/home/amol/Desktop/zmqp/src/github.com/amolb89/cluster/server/serverConfig1.json"
	var sendcount uint64 = 0
	server := new(clust.Serv)
	server.Set(5561,configPath)
	//Server 1
	server1 := new(clust.Serv)
	server1.Set(5562,configPath)

	//Server 2
	server2 := new(clust.Serv)
	server2.Set(5563,configPath)

	//Server 3
	server3 := new(clust.Serv)
	server3.Set(5564,configPath)
	server.Outbox() <- &clust.Envelope{-1,"hello"}
	sendcount = sendcount + 3
	server1.Outbox() <- &clust.Envelope{-1,"hello"}
	sendcount = sendcount + 3	
	server2.Outbox() <- &clust.Envelope{-1,"hello"}
	sendcount = sendcount + 3
	server3.Outbox() <- &clust.Envelope{-1,"hello"}
	sendcount = sendcount + 3	
	var ReceiveCount uint64 = 0
		go func() {
		for {
			<- server.Inbox()
			atomic.AddUint64(&ReceiveCount, 1)
		}
	}()
	go func() {
		for {
			<- server1.Inbox()
			atomic.AddUint64(&ReceiveCount, 1)
		}
	}()

	go func() {
		for {
			<- server2.Inbox()
			atomic.AddUint64(&ReceiveCount, 1)
		}
	}()

	go func() {
		for {
			<- server3.Inbox()
			atomic.AddUint64(&ReceiveCount, 1)
		}
	}()
	time.Sleep(time.Second*5)
	if sendcount != ReceiveCount {
		t.Errorf("Sending and receiving count do not match",sendcount,ReceiveCount)
	}
}



