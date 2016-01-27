package main

import (
	"fmt"
	"log"
	"net"
)

func main(){
	serverAddr := "129.241.187.23:20003"
	
	UDPAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil{
		log.Fatal(err)
	}
	
	conn, err := net.DialUDP("udp", nil, UDPAddr)
	if err != nil{
		log.Fatal(err)
	}
	defer conn.Close()
	
	msg := "pella"
	_, err = conn.Write([]byte(msg + "\x00"))
	if err != nil{
		log.Fatal(err)
	}
	buffer := make([]byte,1024)

	n_bytes, _, _ := conn.ReadFromUDP(buffer)
	fmt.Println(buffer[:n_bytes])
	
}
