/*
 * Copyright 2015 Derek Bever
 *
 * This file is part of ps1rfid.
 *
 * ps1rfid is free software: you can redistribute it and/or modify it under
 * the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
 * for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// The Robotter interface specifies what any bot much implement to
// work with ps1rfid.
type Robotter interface {
	runRobot(shutdown chan bool)
	openDoor()
}

func main() {

	var settingsFile string
	flag.StringVar(&settingsFile, "config", "", "Path to config file. When this is not set it uses default values")
	var testMode bool
	flag.BoolVar(&testMode, "testMode", false, "Use this flag to run this thing in test mode")
	flag.Parse()

	cfg, err := ReadConfig(settingsFile)
	if err != nil {
		fmt.Print(err)
		// TODO what do you actually want to do here? bail out?
		cfg = ConfigDefault
	}

	log.Printf("Test mode flag set to: %v", testMode)
	log.Printf("Config settings:")
	spew.Dump(cfg)

	var thisRobot Robotter
	if testMode {
		var dummy DummyRobot
		thisRobot = dummy
		log.Println("DummyRobot intialized")
	} else {
		var realRobot Robot
		realRobot.configure(cfg)
		thisRobot = realRobot
		log.Println("RealRobot initialized")
	}

	shutdown := make(chan bool, 1)

	//catch SIGINT, SIGKILL for clean shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Caught %v. Shutting down the running goroutine.", sig)
			shutdown <- true
			goto quit
		}
	quit:
		log.Println("Bailin' out of the signal notification goroutine")
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}()

	go thisRobot.runRobot(shutdown)

	serve(thisRobot)
}
