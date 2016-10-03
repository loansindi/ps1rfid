// -*- Mode: Go; indent-tabs-mode: t -*-

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
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config contains information used to connect to resources
type Config struct {
	Version        string
	BoltPath       string `toml:"bolt_path"`
	ServicePort    int    `toml:"service_port"`
	RFIDurl        string `toml:"rfid_url"`
	ToggleDuration int    `toml:"toggle_duration"`
	TogglePin      string `toml:"toggle_pin"`
	SerialName     string `toml:"serial_name"`
	SerialBaud     int    `toml:"serial_baud"`
}

// ConfigDefault holds the default settings
var ConfigDefault = Config{
	BoltPath:       "rfid-tags.db",
	RFIDurl:        "https://members.pumpingstationone.org/rfid/check/FrontDoor",
	ToggleDuration: 5,
	TogglePin:      "P9_11",
	SerialName:     "/dev/ttyUSB0",
	SerialBaud:     9600,
}

// ReadConfig does what it says on the tin.
func ReadConfig(configFile string) (cfg Config, err error) {
	if configFile == "" {
		fmt.Print("Using default config")
		cfg = ConfigDefault
		return
	}
	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		err = fmt.Errorf("Error opening %v: %v", configFile, err)
		return
	}
	err = toml.Unmarshal(contents, &cfg)
	if err != nil {
		err = fmt.Errorf("Error unmarshalling %v: %v", configFile, err)
		return
	}
	return
}
