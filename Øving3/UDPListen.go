package UDP

import (
	"io"
	"log"
	"net"
)

func main(){
	ListenAddr, err := net.ResolveAddr("udp", ":30000")
	if err != nil{
		log.Fatal(err)
	}
	
	buffer	byte[1024]
	listenConn, err := net.ListenUDP("udp", ListenAddr)
	if err != nil{
		log.Fatal(err)
	}
	defer listenConn.Close()
	
	for{
		n_bytes, addr, err := listenConn.ReadFromUDP(buffer)
		fmt.Println("Received: ", string(buf[0:n_bytes], " from ", addr)
		
		if err != nil{
			log.Fatal(err)
		}
	}		
}

