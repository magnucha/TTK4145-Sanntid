package main

import (
	"net"
	"log"
	"time"
	"os/exec"
)

func main(){
	backup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")
	backup.Run()
	
	UDPaddr, err := net.ResolveUDPAddr("udp", "198.168.1.255:9001")
	if err != nil{
		log.Printf(err.Error())
	}
	
	connection, err := net.DialUDP("udp", nil, UDPaddr)
	if err != nil{
		log.Printf("DialUDP error: %s",err.Error())
	}
	
	msg := make([]byte, 1)
	counter := 0
	for{
		log.Printf("%d", counter)
		msg[0] = byte(counter)
		_, err := connection.Write(msg)
		if err != nil{
			log.Printf("UDP write error: %s", err.Error())
		}
		counter++
		time.Sleep(time.Second)
	}
}

