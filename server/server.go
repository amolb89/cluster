package main

import (
	"github.com/amolb89/cluster/clust"
	"strconv"
	 "os"
)

func main(){	
	args := os.Args
	ServerId,_ := strconv.Atoi(args[0])
	configPath := "/home/amol/Desktop/zmqp/src/github.com/amolb89/cluster/server/serverConfig.json"
	server := new(clust.Serv)
	server.Set(ServerId,configPath)
	
}


