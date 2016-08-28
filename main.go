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
	"log"
)

type Robotter interface {
	runRobot()
	openDoor()
}

func main() {

	// var settingsFile string
	// flag.StringVar(&settingsFile, "config", "./config.toml", "Path to the config file")
	// flag.Parse()
	// config, err := cfg.ReadConfig(settingsFile)
	// fmt.Printf("Config: %v", config)

	var testMode bool
	flag.BoolVar(&testMode, "testMode", false, "Use this flag to run this thing in test mode")
	flag.Parse()

	log.Printf("Test mode flag set to: %v", testMode)

	var thisRobot Robotter
	if testMode {
		var dummy DummyRobot
		thisRobot = dummy
		log.Println("DummyRobot intialized")
	} else {
		var realRobot Robot
		realRobot.configure()
		thisRobot = realRobot
		log.Println("RealRobot initialized")
	}

	go thisRobot.runRobot()

	serve(thisRobot)
}
