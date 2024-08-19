package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Hostname : ", hostname)

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, addr := range addrs {
		fmt.Println("IP: ", addr)
	}

	

}