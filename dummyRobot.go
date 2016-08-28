package main

import (
	"log"
	"time"
)

type DummyRobot struct {
}

func (d DummyRobot) configure() {

}

func (d DummyRobot) runRobot() {

	for {
		time.Sleep(10 * time.Second)
	}

}

func (d DummyRobot) openDoor() {
	log.Println("Opening door!")
	time.Sleep(5 * time.Second)
	log.Println("Locking door!")
}
