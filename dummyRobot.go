package main

import (
	"log"
	"time"
)

type DummyRobot struct {
}

func (d DummyRobot) configure() {

}

func (d DummyRobot) runRobot(shutdown chan bool) {

	for {
		select {
		case <-shutdown:
			log.Println("Caught the shutdown signal. Bailing out.")
			goto quit
		default:
			time.Sleep(1 * time.Second)
		}
	}
quit:
	log.Println("Exited the runRobot loop successfully. Later, taters.")
}

func (d DummyRobot) openDoor() {
	log.Println("Opening door!")
	time.Sleep(5 * time.Second)
	log.Println("Locking door!")
}
