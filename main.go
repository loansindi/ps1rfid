package main

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
	//"github.com/pebbe/zmq4"
	"github.com/tarm/goserial"
	"io"
	"net/http"
	"os"
	"time"
)

var cacheDB *bolt.DB

func checkCacheDBForTag(tag string) bool {
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

func addTagToCacheDB(tag string) {
	cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RFIDBucket"))
		err := b.Put([]byte(tag), []byte(tag))
		return err
	})
}

func openDoor(sp gpio.DirectPinDriver) {
	sp.DigitalWrite(1)
	gobot.After(5*time.Second, func() {
		sp.DigitalWrite(0)
	})
}

func main() {

	var code string
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	//NewDirectPinDriver returns a pointer - this wasn't immediately obvious to me
	splate := gpio.NewDirectPinDriver(beagleboneAdaptor, "splate", "P9_13")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	u, err := serial.OpenPort(c)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// the anonymous function here allows us to call openDoor with splate remaining in scope
	go http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
		openDoor(*splate)
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

		// Now open the cache db to check if it's already here
		cacheDB, err = bolt.Open("rfid-tags.db", 0600, nil)
		if err != nil {
			fmt.Println(err)
		}

		cacheDB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("RFIDBucket"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})

		// Before checking the site for the code, let's check our cache
		if checkCacheDBForTag(code) == false {
			var request bytes.Buffer
			request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
			request.WriteString(code)
			resp, err := http.Get(request.String())
			if err != nil {
				fmt.Printf("Whoops!")
				os.Exit(1)
			}
			if resp.StatusCode == 200 {

				// We got 200 back, so we're good to add this
				// tag to the cache
				addTagToCacheDB(code)

				fmt.Println("Success!")
				code = ""
				openDoor(*splate)
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
			openDoor(*splate)
		}

		cacheDB.Close()
	}
}
