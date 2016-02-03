package main

import (
	"log"
	"network"
	"config"
)

func main(){
	ch_incoming_msg := make(chan config.Message)
	ch_outgoing_msg := make(chan config.Message)
	network.Network_Init(ch_incoming_msg, ch_outgoing_msg)
}
