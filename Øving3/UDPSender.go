package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main(){
	serverAddr := "129.241.187.23:20010"
	
	//Set up send socket
	remoteAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil{
		log.Fatal(err)
	}
	socketSend, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil{
		log.Fatal(err)
	}
	
	//Set up receive socket
	port, err := net.ResolveUDPAddr("udp", ":20010")
	if err != nil{
		log.Fatal(err)
	}
	socketReceive, err := net.ListenUDP("udp", port)
	if err != nil{
		log.Fatal(err)
	}



	defer socketReceive.Close()
	defer socketSend.Close()
	
	msg := "smmella"
	_, err = socketSend.Write([]byte(msg + "\x00"))
	if err != nil{
		log.Fatal(err)
	}
	buffer := make([]byte,1024)
	for{
		n_bytes, _, _ := socketReceive.ReadFromUDP(buffer)
		fmt.Println(string(buffer[:n_bytes]))
		
		msg := "smmella"
		_, err = socketSend.Write([]byte(msg + "\x00"))

		time.Sleep(time.Second)
	}
}
