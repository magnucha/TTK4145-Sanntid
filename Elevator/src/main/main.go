package main

import (
	"config"
	"fsm"
	"hardware"
	"log"
	"network"
	"queue"
	"time"
	"os"
	"os/exec"
)

var ch_incoming_msg = make(chan config.Message)
var ch_outgoing_msg = make(chan config.Message)
var ch_new_order = make(chan config.ButtonStruct)
//var ch_del_order = make(chan config.Message)
var ch_button_pressed = make(chan config.ButtonStruct)
var ch_floor_poll = make(chan int)
var ch_new_elev = make(chan string)
var ch_main_alive = make(chan bool)

func main() {
	if _,err := os.Open(config.QUEUE_FILENAME); err == nil {
		Backup_Hold()
		queue.File_Read(config.QUEUE_FILENAME)
		network.Init(ch_outgoing_msg, ch_incoming_msg, ch_new_elev, ch_main_alive)
	} else {
		network.Init(ch_outgoing_msg, ch_incoming_msg, ch_new_elev, ch_main_alive)
		time.Sleep(time.Millisecond)
		if _,err := os.Create(config.QUEUE_FILENAME); err != nil {
			log.Fatal("FATAL: Could not create queue file!")
		}
	}
	backup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main/main.go")
	backup.Run()

	if !hardware.Elev_Init() {
		log.Fatal("Unable to initialize elevator hardware!")
	}
	hardware.Elev_Set_Motor_Direction(fsm.Choose_New_Direction())
	go Message_Server()
	go Channel_Server()
	go hardware.Read_Buttons(ch_button_pressed)
	go hardware.Set_Lights()
	go hardware.Floor_Poller(ch_floor_poll)
	go State_Spammer()
	go fsm.Event_Order_Received()
	fsm.Init(ch_outgoing_msg, ch_new_order)

	log.Printf("Elev addr: %s", config.Laddr)

	for {
		time.Sleep(3 * time.Second)
	}

}

func Message_Server() {
	for {
		msg := <-ch_incoming_msg
		switch msg.Msg_type {
		//case config.ACK:
		//	Increment_Ack_Counter(msg)	//Not yet implemented
		case config.STATE_UPDATE:
			already_active := false
			for addr, _ := range config.Active_elevs {
				if msg.Raddr == addr {
					already_active = true
					break
				}
			}
			if !already_active {
				ch_new_elev <- msg.Raddr
				time.Sleep(10 * time.Microsecond)
			}
			State_Copy(config.Active_elevs[msg.Raddr], &msg.State)
			config.Active_elevs[msg.Raddr].Timer.Reset(config.TIMEOUT_REMOTE)
		case config.ADD_ORDER:
			ch_new_order <- msg.Button
		case config.DELETE_ORDER:
			log.Println("Remote delete order received")
			queue.Delete_Order(msg.Button.Floor, ch_outgoing_msg, false)
		}
	}
}

func Channel_Server() {
	for {
		select {
		case button := <-ch_button_pressed:
			ch_new_order <- button
			if button.Button_type != config.BUTTON_COMMAND {
				ch_outgoing_msg <- config.Message{Msg_type: config.ADD_ORDER, Button: button}
			}
		case floor := <-ch_floor_poll:
			fsm.Event_Reached_Floor(floor)
		case raddr := <-ch_new_elev:
			Set_Active(raddr)
		}
	}
}

func State_Spammer() {
	for {
		time.Sleep(500 * time.Millisecond)
		ch_outgoing_msg <- config.Message{Msg_type: config.STATE_UPDATE, State: *config.Active_elevs[config.Laddr]}
	}
}

func Set_Active(raddr string) {
	for addr, _ := range config.Active_elevs {
		if addr == raddr {
			return
		}
	}
	killer := func() {
		delete(config.Active_elevs, raddr)
		queue.Reassign_Orders(raddr, ch_new_order)
	}
	config.Active_elevs[raddr] = &config.ElevState{Is_idle: true, Door_open: false, Direction: config.DIR_STOP, Last_floor: -1, Timer: time.AfterFunc(config.TIMEOUT_REMOTE, killer)}
}

func State_Copy(a *config.ElevState, b *config.ElevState) {
	a.Is_idle = b.Is_idle
	a.Door_open = b.Door_open
	a.Direction = b.Direction
	a.Last_floor = b.Last_floor
}

func Backup_Hold() {
	var ch_reset = make(chan config.NetworkMessage)
	timer_alive := time.NewTimer(config.TIMEOUT_LOCAL)
	conn := network.UDP_Create_Listen_Socket(config.UDP_ALIVE_PORT)
	defer conn.Close()
	go network.UDP_Receive(conn, ch_reset)
	
	for {
		select {
		case msg := <- ch_reset:
			if string(msg.Data[:len(config.UDP_PRESENCE_MSG)]) == config.UDP_PRESENCE_MSG {
				timer_alive.Reset(config.TIMEOUT_LOCAL)
			}
		case <- timer_alive.C:
			return
		}
		time.Sleep(100*time.Millisecond)
	}
}

/*
Feiltoleranse:
	- Cosmic rays
		- DAFUQ???
		- acceptance test on variables

*/
