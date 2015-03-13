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
		log.Fatal(err)
	}
	r = client
}

func main() {
	// pull an INFO all call into a RedisInfoAll struct
	all, err := r.Info()
	if err != nil {
		log.Fatal("unable to connect and get info")
	}
	// Accessing the Server section is done via all.Server
	fmt.Printf("Redis Server Version: %s\n", all.Server.Version)
	fmt.Printf("Redis Server Role: %s\n", all.Replication.Role)
	// Accessing the Replication section is done via all.Replication
	switch all.Replication.Role {
	case "slave":
		fmt.Printf("Slave of: %s:%d\n", all.Replication.MasterHost, all.Replication.MasterPort)
	case "master":
		println("Is a Master")
	}
	if all.Replication.ConnectedSlaves > 0 {
		fmt.Println("Slaves")
		fmt.Println("======")
		fmt.Printf("Number Connected: %d\n", all.Replication.ConnectedSlaves)
		for _, slave := range all.Replication.Slaves {
			// slave here is an InfoSlaves instance
			println(" Slave")
			fmt.Printf("    IP: %s\n", slave.IP)
			fmt.Printf("    Port: %d\n", slave.Port)
			fmt.Printf("    State: %s\n", slave.State)
			fmt.Printf("    Lag: %d\n", slave.Lag)
			fmt.Printf("    Offset: %d\n", slave.Offset)
		}
	}
	// Accessing the Memory section is done via all.Memory
	fmt.Printf("Redis Used Memory: %db\n", all.Memory.UsedMemory)
	fmt.Printf("Redis Used Memory Human: %s\n", all.Memory.UsedMemoryHuman)

}
