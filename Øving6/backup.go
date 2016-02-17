package main

import (
	"log"
	"encoding/binary"
	"time"
	"net"
	"os/exec"
)

func listen(UDP *net.UDPConn, ch chan<- int) {
	buf := make([]byte, 1024)
	for {
		UDP.ReadFromUDP(buf[:])
		
		rec,_ := binary.Uvarint(buf)
		ch <- int(rec)
		log.Printf("Recieved %d", int(rec))
		time.Sleep(100*time.Millisecond)
	}
}

func backup(UDP *net.UDPConn) int {
	lastVal := 0
	ch := make(chan int)
	go listen(UDP, ch)
	for {
		select {
			case lastVal = <- ch:
				break;
			case <-time.After(3*time.Second):
				return lastVal
		}
	}
}

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":20003")
	listen, _ := net.ListenUDP("udp", addr)
	counter := backup(listen) + 1 //Offset to avoid transmitting the same number twice
	listen.Close()
	
	addr, _ = net.ResolveUDPAddr("udp","129.241.187.150:20003")
	broadcast, _ := net.DialUDP("udp", nil, addr)
	
	backup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")
	backup.Run()
	msg := make([]byte, 1)
	for{
		log.Printf("%d", counter)
		msg[0] = byte(counter)
		_, err := broadcast.Write(msg)
		if err != nil{
			log.Printf("UDP write error: %s", err.Error())
		}
		counter++
		time.Sleep(time.Second)
	}
}
