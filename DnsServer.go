package main

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)


var ctx = context.Background()


func HandleDNSquery(request []byte, upstreamDNS string , rdb *redis.Client) ([]byte, error) {

	var query DNSQuery
	err := query.DNSMessageFromBytes(request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DNS question: %v", err)
	}

    for i := 0 ; i < int(query.Header.QDCount) ; i++ {
        fmt.Printf("query domain number %d : %s ,", i , query.Questions[0].domain)
    }

	cachedResponse, err := getCachedDNSResponse(ctx , rdb, query.Questions[0].domain)
		if err == nil && len(cachedResponse) > 0 {
			return cachedResponse , nil
	}
	
	upstreamAddr, err := net.ResolveUDPAddr("udp", upstreamDNS)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve upstream DNS address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, upstreamAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to upstream DNS server: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send DNS request to upstream server: %v", err)
	}

	response := make([]byte, 512)
	n, _, err := conn.ReadFromUDP(response)
	if err != nil {
		return nil, fmt.Errorf("failed to receive DNS response from upstream server: %v", err)
	}
	response = response[:n]

    ttl , err := ExtractTTL(response)
	if err == nil && ttl > 0 {
		err = cacheDNSResponse(ctx , rdb, query.Questions[0].domain, response, ttl)
		if err != nil {
			fmt.Println("error caching DNS response:", err)
		}else{
		fmt.Printf("Cached response for %s with TTL %v seconds\n", query.Questions[0].domain, ttl.Seconds())
		}
	}

	return response, nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
		DB:   0,                
	})

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
		fmt.Printf("recieved msg from %s : %s 	\n" , senderAddr , string(buffer[:n]))
        upsteam , err := GetUpstream()
        if err != nil{
            fmt.Printf("the world is down")
        }
        responce , err := HandleDNSquery(buffer[:n] , upsteam , rdb)
        if err != nil {
            fmt.Printf("the world is down vol 2")
        }



		_ , err = conn.WriteToUDP(responce , senderAddr)
		if err != nil {
			fmt.Println("hello negro didnt sent")
		}
	}
}