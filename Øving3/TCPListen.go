package TCP

import (
	"net"
	"bufio"
	"os"
	"fmt"
)

func TCP_receive(conn net.Conn) string {
	//Wait for message ending in '\0'
	msg, _ := bufio.NewReader(conn).ReadString('\0')
	return msg
}

func TCP_send(conn net.Conn, msg string) {
	conn.Write([]byte(msg + '\0'))
}

func main() {
	fmt.Println("Launching TCP server...")
	input := bufio.NewReader(os.Stdin)
	
	//Listen on port 33546
	listener,_ := net.Listen("tcp", ":33546")
	
	//Set to accept connections
	conn, _ := listener.Accept()
	
	for {
		//Send inititial message
		fmt.Print("Enter message to send to server: ")
		msg_send, _ := input.ReadString('\n')
		TCP_send(conn, msg_send)
		
		msg_rcpt := TCP_receive(conn)
		
		//DEBUG: Print received message to console
		fmt.Print("Message received: ", string(msg_rcpt))
		
		//Prevent spamming the network
		sleep(1)
	}
}

