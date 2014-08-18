package main

import (
	"fmt"
	"log"
	"time"

	rclient "github.com/therealbill/libredis/client"
)

var (
	network  = "tcp"
	address  = "127.0.0.1:6379"
	db       = 1
	password = ""
	timeout  = 5 * time.Second
	maxidle  = 1
	r        *rclient.Redis
)

func init() {
	client, err := rclient.Dial("127.0.0.1", 6379)
	if err != nil {
		panic(err)
	}
	r = client
}

func main() {
	all, err := r.Info()
	if err != nil {
		log.Fatal("unable to connect and get info")
	}
	fmt.Printf("Redis Server Version: %s\n", all.Server.Version)
	fmt.Printf("Redis Server Role: %s\n", all.Replication.Role)
	if all.Replication.ConnectedSlaves > 0 {
		fmt.Println("Slaves:")
		fmt.Printf("\tNumber Connected: %d\n", all.Replication.ConnectedSlaves)
		for _, slave := range all.Replication.Slaves {
			fmt.Printf("\tSlave: %+v\n", slave)
		}
	}
	if all.Replication.Role == "slave" {
		fmt.Printf("Master: %s:%d\n", all.Replication.MasterHost, all.Replication.MasterPort)
	}

	fmt.Printf("Redis Used Memory: %db\n", all.Memory.UsedMemory)
	fmt.Printf("Redis Used Memory Human: %s\n", all.Memory.UsedMemoryHuman)

}
