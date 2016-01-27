package main

import (
	"log"
	"net"
	"fmt"
)

func main(){
	ListenAddr, err := net.ResolveUDPAddr("udp", ":30000")
	if err != nil{
		log.Fatal(err)
	}
	
	buffer := make([]byte,1024)
	listenConn, err := net.ListenUDP("udp", ListenAddr)
	if err != nil{
		log.Fatal(err)
	}
	defer listenConn.Close()
	
	//for{
		n_bytes, addr, err := listenConn.ReadFromUDP(buffer)
		fmt.Println("Received: ", string(buffer[0:n_bytes]), " from ", addr)
		
		if err != nil{
			log.Fatal(err)
		}
	//}
}
