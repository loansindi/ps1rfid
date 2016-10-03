package main

import (
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

// Robot contains resources used to interact with the BeagleBone and Cache.
type Robot struct {
	cacheDB     *bolt.DB
	strikePlate *gpio.DirectPinDriver
	serialPort  io.ReadWriter
	cfg         Config
}

func (r *Robot) configure(cfg Config) {
	r.cfg = cfg
	// Now open the cache db to check if it's already here
	var boltErr error
	r.cacheDB, boltErr = bolt.Open(cfg.BoltPath, 0600, nil)
	if boltErr != nil {
		log.Fatalf("Unable to open the cacheDB: %+v", boltErr)
	}

	//Now we initialize the RFIDBucket bucket if it doesn't already exist
	r.cacheDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("RFIDBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	//Configure the DirectPinDriver
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	//NewDirectPinDriver returns a pointer - this wasn't immediately obvious to me
	r.strikePlate = gpio.NewDirectPinDriver(beagleboneAdaptor, "splate", cfg.TogglePin)

	//Configure the serial port
	c := &serial.Config{Name: cfg.SerialName, Baud: cfg.SerialBaud}
	var serialErr error
	r.serialPort, serialErr = serial.OpenPort(c)
	if serialErr != nil {
		log.Fatalf("Unalbe to open serial port: %+v", serialErr)
	}
}

func (r Robot) openDoor() {
	r.strikePlate.DigitalWrite(1)
	duration := time.Duration(r.cfg.ToggleDuration) * time.Second
	gobot.After(duration, func() {
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

func (r Robot) runRobot(shutdown chan bool) {
	defer r.cacheDB.Close()
	for {
		// Check to see if we need to exit this loop
		select {
		case <-shutdown:
			// break out of the for loop, which will trigger the r.cacheDB.Close()
			log.Println("Caught the shutdown signal. Bailing out.")
			goto quit
		default: // if we don't have anything on shutdown, keep on keepin' on.
			code := r.serialRead()
			// Before checking the site for the code, let's check our cache
			if r.checkCacheDBForTag(code) {
				log.Printf("%s scanned in via the cache successfully.", code)
				r.openDoor()
			} else if r.checkPS1ForTag(code) {
				log.Printf("%s scanned in via members.ps1.org successfully", code)
				r.openDoor()
			} else {
				log.Printf("%s was found in neither the cache not the ps1 member site.", code)
			}
		}
	}
quit:
	log.Println("Exited the run loop. Later taters.")
}

func (r Robot) checkPS1ForTag(code string) bool {
	rfidCheckURL := fmt.Sprintf("%s/%s", r.cfg.RFIDurl, code)
	resp, err := http.Get(rfidCheckURL)
	if err != nil {
		log.Printf("Unable to access %s for this reason: %+v", rfidCheckURL, err)
		return false
	}
	if resp.StatusCode == http.StatusOK {
		// We got 200 back, so we're good to add this
		// tag to the cache
		r.addTagToCacheDB(code)

		log.Printf("%s found in the database and added to cache.", code)
		return true
	}
	if resp.StatusCode == http.StatusForbidden {
		log.Printf("%s tried to scan in, but mebership was expired.", code)
		return false
	}
	log.Printf("%s tried to scan in but encountered status %s", code, http.StatusText(resp.StatusCode))
	return false
}

func (r Robot) checkCacheDBForTag(tag string) bool {
	val := ""
	r.cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RFIDBucket"))
		val = string(b.Get([]byte(tag)))
		return nil
	})

	if val != "" {
		return true
	}

	return false
}

func (r Robot) addTagToCacheDB(tag string) {
	r.cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RFIDBucket"))
		err := b.Put([]byte(tag), []byte(tag))
		return err
	})
}
