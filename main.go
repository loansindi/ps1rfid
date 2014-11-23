package main

import (
	"io"
	"time"
	"bytes"
	"fmt"
	"github.com/tarm/goserial"
	"net/http"
	"os"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"

)

var code string

func main() {
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	splate := gpio.NewDirectPinDriver(beagleboneAdaptor, "splate", "P9_11")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	u, err := serial.OpenPort(c)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	go http.HandleFunc("/", displayCode)
//	go http.HandleFunc("/opensesame", openDoor)
	go http.ListenAndServe(":8080", nil)
	buf := make([]byte, 16)
	for {
		n, err := io.ReadFull(u, buf)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		code = string(buf[1 : n-3])
		var request bytes.Buffer
		request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
		request.WriteString(code)
		resp, err := http.Get(request.String())
		if err != nil {
			fmt.Printf("Whoops!")
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			fmt.Println("Success!")
			code = ""
			splate.DigitalWrite(1)
			gobot.After(5*time.Second, func() {
				splate.DigitalWrite(0)
			})
		} else if resp.StatusCode == 403 {
			fmt.Println("Membership status: Expired")
		} else {
			fmt.Println("Code not found")
		}
	}

}
/*
func openDoor(w http.ResponseWriter, r *http.Request) {
	splate.DigitalWrite(1)
	gobot.After(5*time.Second, func() {
		splate.DigitalWrite(0)
})

}
*/

func displayCode(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Code: "))
	w.Write([]byte(code))
}

func openDoor() {

	

}
/*
func bot() {

//	gbot := gobot.NewGobot()
	
	work := func() {
		splate.DigitalWrite(1)
		gobot.After(1*time.Second, func() {
			splate.DigitalWrite(0)
		})

}}
*/
