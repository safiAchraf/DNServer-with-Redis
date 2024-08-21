package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	addr , err := net.ResolveUDPAddr("udp" , ":8080")
	if err != nil {
		fmt.Println("error creating the addr")
	}
	conn , err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Println("error listening")
	}

	buffer := make([]byte , 1024)


	for {
		n , senderAddr , err := conn.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("error recieve packet")
			continue
		}
		fmt.Println("recieved msg from %s : %s \n" , senderAddr , string(buffer[:n]))

		message := []byte("hellow negro")
		n , err = conn.WriteTo(message , senderAddr)
		if err != nil {
			fmt.Println("hello negreo didnt sent")
		}
	}
}