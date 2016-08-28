package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/tarm/goserial"
)

type Robot struct {
	cacheDB     *bolt.DB
	strikePlate *gpio.DirectPinDriver
	serialPort  io.ReadWriter
}

func (r *Robot) configure() {
	// Now open the cache db to check if it's already here
	var boltErr error
	r.cacheDB, boltErr = bolt.Open("rfid-tags.db", 0600, nil)
	if boltErr != nil {
		log.Fatalf("Unable to open the cacheDB: %+v", boltErr)
	}

	//Configure the DirectPinDriver
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	//NewDirectPinDriver returns a pointer - this wasn't immediately obvious to me
	r.strikePlate = gpio.NewDirectPinDriver(beagleboneAdaptor, "splate", "P9_11")

	//Configure the serial port
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	var serialErr error
	r.serialPort, serialErr = serial.OpenPort(c)
	if serialErr != nil {
		log.Fatalf("Unalbe to open serial port: %+v", serialErr)
	}
}

func (r Robot) openDoor() {
	r.strikePlate.DigitalWrite(1)
	gobot.After(5*time.Second, func() {
		r.strikePlate.DigitalWrite(0)
	})
}

func (r Robot) serialRead() string {
	buf := make([]byte, 16)
	n, err := io.ReadFull(r.serialPort, buf)
	if err != nil {
		log.Fatalf("Unable to read bytes from seral port: %+v", err)
	}
	// We need to strip the stop and start bytes from the tag, so we only assign a certain range of the slice
	return string(buf[1 : n-3])
}

func (r Robot) runRobot() {
	defer r.cacheDB.Close()
	for {

		code := r.serialRead()

		// I'm not 100% sure what this code is doing. Is it initalizing `RFIDBucket` if it doesn't exist?
		r.cacheDB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("RFIDBucket"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})

		// Before checking the site for the code, let's check our cache
		if checkCacheDBForTag(code, r.cacheDB) == false {
			var request bytes.Buffer
			request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
			request.WriteString(code)
			resp, err := http.Get(request.String())
			if err != nil {
				log.Fatalf("Unable to reach https://members.pumpingstationone.org/rfid/check/FrontDoor/ : %+v", err)
			}
			if resp.StatusCode == 200 {

				// We got 200 back, so we're good to add this
				// tag to the cache
				addTagToCacheDB(code, r.cacheDB)

				fmt.Println("Success!")
				code = ""
				r.openDoor()
			} else if resp.StatusCode == 403 {
				fmt.Println("Membership status: Expired")
			} else {
				fmt.Println("Code not found")
			}
		} else {
			// If we're here, we found the tag in the cache, so
			// let's just go and open the door for 'em
			fmt.Println("Success!")
			code = ""
			r.openDoor()
		}

	}
}

func checkCacheDBForTag(tag string, cacheDB *bolt.DB) bool {
	val := ""
	cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RFIDBucket"))
		val = string(b.Get([]byte(tag)))
		return nil
	})

	if val != "" {
		return true
	}

	return false
}

func addTagToCacheDB(tag string, cacheDB *bolt.DB) {
	cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RFIDBucket"))
		err := b.Put([]byte(tag), []byte(tag))
		return err
	})
}
