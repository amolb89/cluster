package clust

import (
	zmq "github.com/pebbe/zmq4"
	"encoding/json"	
	"strconv"
)

const (BROADCAST = -1)

type Envelope struct {
	Pid int
	Msg string
}

type Server interface {
	Pid() int
	Peers() map[int]string
	Outbox() chan *Envelope
	Inbox() chan *Envelope
}

type Serv struct{
	pid int
	peers map[int]string
	outbox chan *Envelope
	inbox chan *Envelope
}

func (s *Serv) Set(id int, configPath string) {
	cfig := new(Config)
	//Load configuration from config file
	LoadConfig(configPath,cfig)
	Peers := make(map[int]string)
	
	//Initialize list of peers
	for servers,_ := range cfig.Servers {
		key,_ := strconv.Atoi(servers)
		Peers[key]=cfig.Servers[servers]
	}

	//Create inbox and outbox channels
	s.outbox = make(chan *Envelope,100)
	s.inbox = make(chan *Envelope,100)
	s.pid = id
	s.peers = Peers

	//Start receiving and sending goroutines
	go s.Proc_recv()
	go s.Proc_send()
}

func (s *Serv) Pid() int {
	return s.pid
}

func (s *Serv) Peers() map[int]string {
	return s.peers
}

func (s *Serv) Outbox() chan *Envelope {
	return s.outbox
}

func (s *Serv) Inbox() chan *Envelope {
	return s.inbox
}

func (s *Serv) Proc_send() {
	for {
		msg := <- s.Outbox()
		if msg.Pid == -1 {
			for id, _ := range s.Peers() {
				if id != s.Pid() {
					requester, _ := zmq.NewSocket(zmq.REQ)
					servAddr := "tcp://"+s.Peers()[id]
					requester.Connect(servAddr)
					msg.Pid = s.Pid()
					msgBytes,_ := json.Marshal(msg)
					requester.SendBytes(msgBytes, 0)
					requester.Recv(0)
					requester.Close()
				}
			}
		} else {
			requester, _ := zmq.NewSocket(zmq.REQ)
			servAddr := "tcp://"+s.Peers()[msg.Pid]
			requester.Connect(servAddr)
			msg.Pid = s.Pid()
			msgBytes,_ := json.Marshal(msg)
			requester.SendBytes(msgBytes, 0)
			requester.Recv(0)
			requester.Close()
		}
	}
}

func (s *Serv) Proc_recv() {
        responder, _ := zmq.NewSocket(zmq.REP)
        defer responder.Close()
	recAddr := "tcp://"+s.Peers()[s.Pid()]
        responder.Bind(recAddr)
        for {
                msgBytes, _ := responder.RecvBytes(0)
		var msg Envelope
		json.Unmarshal(msgBytes, &msg)
		s.Inbox() <- &msg
                // Send reply back to client
                reply := "1"
                responder.Send(reply, 0)
        }
}
