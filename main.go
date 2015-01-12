package main

import (
	"bytes"
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/tarm/goserial"
	"io"
	"net/http"
	"os"
	"time"
	zmq "github.com/pebbe/zmq4"
)

func openDoor(sp gpio.DirectPinDriver, publisher *zmq.Socket) {
	sp.DigitalWrite(1)
	publisher.SendMessage("door.state.unlock", "Door Unlocked")
	gobot.After(5*time.Second, func() {
		sp.DigitalWrite(0)
		publisher.SendMessage("door.state.lock", "Door Locked")
	})
}

func main() {
	var code string
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	//NewDirectPinDriver returns a pointer - this wasn't immediately obvious to me
	splate := gpio.NewDirectPinDriver(beagleboneAdaptor, "splate", "P9_11")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	u, err := serial.OpenPort(c)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	//Configure ZMQ publisher
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	publisher.Bind("tcp://*:5556")
	//Configure ZMQ publisher
	go http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Code: "))
		w.Write([]byte(code))
	})
	// the anonymous function here allows us to call openDoor with splate remaining in scope
	go http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
		openDoor(*splate, publisher)
	})
	go http.ListenAndServe(":8080", nil)
	buf := make([]byte, 16)
	for {
		n, err := io.ReadFull(u, buf)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		// We need to strip the stop and start bytes from the tag, so we only assign a certain range of the slice
		code = string(buf[1 : n-3])
		var request bytes.Buffer
		request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
		request.WriteString(code)
		resp, err := http.Get(request.String())
		if err != nil {
			fmt.Printf("Whoops!")
			publisher.SendMessage("door.rfid.error", fmt.Sprintf("Auth Server Error: %s", err))
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			fmt.Println("Success!")
			publisher.SendMessage("door.rfid.accept", "RFID Accepted")
			code = ""
			openDoor(*splate, publisher)
		} else if resp.StatusCode == 403 {
			fmt.Println("Membership status: Expired")
			publisher.SendMessage("door.rfid.deny", "RFID Denied")
		} else {
			fmt.Println("Code not found")
			publisher.SendMessage("door.rfid.deny", "RFID not found")
		}
	}

}
