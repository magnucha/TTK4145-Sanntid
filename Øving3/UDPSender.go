package UDP

import (
	"io"
	"log"
	"net"
)

func main(){
	serverAddr string = ":20003"
	
	senderAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil{
		log.Fatal(err)
	}
	
	conn, err := net.DialUDP("udp", nil, senderAddr)
	if err != nil{
		log.Fatal(err)
	}
	defer conn.Close()
	
	msg := "pella"
	_, err := conn.WriteToUDP(msg, serverAddr)
	if err != nil{
		fmt.Println(err)
	}
	
}
