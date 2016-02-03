package config






type NetworkMessage struct {
	raddr string 	//The remote address we are receiving from, on form IP:port. 
	data []byte			
	length int			//Length of received data, don't care when transmitting
}

TCP_connections := make([]conn.TCPConn,1)
