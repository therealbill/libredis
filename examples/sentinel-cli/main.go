package main

import (
	"fmt"
	"log"
	"time"

	rclient "github.com/therealbill/libredis/client"
)

var (
	network  = "tcp"
	address  = "127.0.0.1:26379"
	db       = 1
	password = ""
	timeout  = 5 * time.Second
	maxidle  = 1
	r        *rclient.Redis
)

func init() {
	client, err := rclient.DialWithConfig(&rclient.DialConfig{Address: address})
	if err != nil {
		log.Fatal("Unable to connect to Sentinel instance:", err)
	}
	r = client
}

func main() {
	// pull an INFO all call into a RedisInfoAll struct
	all, err := r.SentinelInfo()
	if err != nil {
		log.Fatal("unable to connect and get info:", err)
		return
	}
	// Accessing the Server section is done via all.Server
	fmt.Printf("Redis Server Version: %s\n", all.Server.Version)
	fmt.Printf("Redis Server Mode: %s\n", all.Server.Mode)
	if all.Server.Mode != "sentinel" {
		log.Fatal("Node is NOT a sentinel instance, aborting.")
	}
	// To get the list of managers (pods) under management, call
	// SentinelMasters. It returns a MasterInfo struct.
	pods, err := r.SentinelMasters()
	if err != nil {
		log.Fatal("Unable to run SENTINEL MASTERS command;", err)
	}
	fmt.Printf("Managed Pod count: %d\n", len(pods))
	for _, pod := range pods {
		fmt.Println("Pod Name:", pod.Name)
		fmt.Println("Pod IP:", pod.IP)
		fmt.Println("Pod Port:", pod.Port)
		fmt.Println("Pod Slave Count:", pod.NumSlaves)
		// We can easily to testing of the conditions reported
		// Here we see if our master has any slaves connected This could be
		// extended to talk to connected slaves to get their slave-priority.
		// This would allow us to validate we have promotable slaves.
		if pod.NumSlaves == 0 {
			fmt.Println("!!WARNING!!\n\tThis pod has no slaves. Failover is not possible!")
		}
		// Here we see if our sentinel constellation has quorum on this pod
		fmt.Println("Pod Quorum:", pod.Quorum)
		if pod.Quorum <= pod.NumOtherSentinels {
			fmt.Println("Quorum is possible")
		} else {
			fmt.Printf("!!CRITICAL!!\n\tQuorum is NOT possible! Need %d other sentinels, have %d\n", pod.Quorum, pod.NumOtherSentinels)
		}
		fmt.Println()
	}

}
