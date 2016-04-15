package main

import (
	"config"
	"fsm"
	"hardware"
	"log"
	"network"
	"os"
	"os/exec"
	"queue"
	"time"
)

var ch_incoming_msg = make(chan config.Message)
var ch_outgoing_msg = make(chan config.Message)
var ch_new_order = make(chan config.ButtonStruct)
var ch_button_pressed = make(chan config.ButtonStruct)
var ch_floor_poll = make(chan int)
var ch_new_elev = make(chan string)
var ch_main_alive = make(chan bool)

func main() {
	if _, err := os.Open(config.QUEUE_FILENAME); err == nil {
		BackupHold()
		queue.FileRead(config.QUEUE_FILENAME)
		network.Init(ch_outgoing_msg, ch_incoming_msg, ch_new_elev, ch_main_alive)
	} else {
		network.Init(ch_outgoing_msg, ch_incoming_msg, ch_new_elev, ch_main_alive)
		time.Sleep(time.Millisecond)
		if _, err := os.Create(config.QUEUE_FILENAME); err != nil {
			log.Fatal("FATAL: Could not create queue file!")
		}
	}
	backup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main/main.go")
	backup.Run()

	if !hardware.Init() {
		log.Fatal("Unable to initialize elevator hardware!")
	}
	hardware.SetMotorDirection(fsm.ChooseNewDirection())
	go MessageServer()
	go hardware.ReadButtons(ch_button_pressed)
	go hardware.SetLights()
	go hardware.FloorPoller(ch_floor_poll)
	go StateSpammer()
	go fsm.EventOrderReceived()
	fsm.Init(ch_outgoing_msg, ch_new_order)

	log.Printf("Elev addr: %s", config.Laddr)

	for {
		select {
		case button := <-ch_button_pressed:
			ch_new_order <- button
			if button.Button_type != config.BUTTON_COMMAND {
				ch_outgoing_msg <- config.Message{Msg_type: config.ADD_ORDER, Button: button}
			}
		case floor := <-ch_floor_poll:
			fsm.EventReachedFloor(floor)
		}
	}
}

func MessageServer() {
	for {
		msg := <-ch_incoming_msg
		switch msg.Msg_type {
		case config.STATE_UPDATE:
			already_active := false
			for addr, _ := range config.Active_elevs {
				if msg.Raddr == addr {
					already_active = true
					break
				}
			}
			if !already_active {
				SetActive(msg.Raddr)
			}
			StateCopy(config.Active_elevs[msg.Raddr], &msg.State)
			config.Active_elevs[msg.Raddr].Timer.Reset(config.TIMEOUT_REMOTE)
		case config.ADD_ORDER:
			ch_new_order <- msg.Button
			ACK_msg := msg
			ACK_msg.Msg_type = config.ACK
			ch_outgoing_msg <- ACK_msg
		case config.DELETE_ORDER:
			queue.DeleteOrder(msg.Button.Floor, ch_outgoing_msg, false)
			ACK_msg := msg
			ACK_msg.Msg_type = config.ACK
			ch_outgoing_msg <- ACK_msg
		}
	}
}

func StateSpammer() {
	for {
		time.Sleep(150 * time.Millisecond)
		ch_outgoing_msg <- config.Message{Msg_type: config.STATE_UPDATE, State: *config.Active_elevs[config.Laddr]}
	}
}

func SetActive(raddr string) {
	for addr, _ := range config.Active_elevs {
		if addr == raddr {
			return
		}
	}
	elev_killer := func() {
		delete(config.Active_elevs, raddr)
		queue.ReassignOrders(raddr, ch_new_order)
	}
	config.Active_elevs[raddr] = &config.ElevState{Is_idle: true, Door_open: false, Direction: config.DIR_STOP, Last_floor: -1, Timer: time.AfterFunc(config.TIMEOUT_REMOTE, elev_killer)}
}

func StateCopy(a *config.ElevState, b *config.ElevState) {
	a.Is_idle = b.Is_idle
	a.Door_open = b.Door_open
	a.Direction = b.Direction
	a.Last_floor = b.Last_floor
}

func BackupHold() {
	var ch_reset = make(chan config.NetworkMessage)
	timer_alive := time.NewTimer(config.TIMEOUT_LOCAL)
	conn := network.UDPCreateListenSocket(config.UDP_ALIVE_PORT)
	defer conn.Close()
	go network.UDPReceive(conn, ch_reset)

	for {
		select {
		case msg := <-ch_reset:
			if string(msg.Data[:len(config.UDP_BACKUP_MSG)]) == config.UDP_BACKUP_MSG {
				timer_alive.Reset(config.TIMEOUT_LOCAL)
			}
			if string(msg.Data[:11]) == "Kill backup" {
				log.Fatal("Kill command recieved!")
			}
		case <-timer_alive.C:
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}