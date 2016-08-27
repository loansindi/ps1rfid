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

package ps1rfid

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Version        string
	ServicePort    int    `toml:"service_port"`
	RFIDurl        string `toml:"rfid_url"`
	RFIDRresource  string `toml:"rfid_resource"`
	ToggleDuration int    `toml:"toggle_duration"`
	TogglePin      string `toml:"toggle_pin"`
	SerialName     string `toml:"serial_name"`
	SerialBaud     int    `toml:"serial_baud"`
}

func ReadConfig(configFile string) (Config, error) {
	var config Config
	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, fmt.Errorf("Error opening %v: %v", configFile, err)
	}
	err = toml.Unmarshal(contents, &config)
	if err != nil {
		return config, fmt.Errorf("Error unmarshalling %v: %v", configFile, err)
	}
	return config, nil
}
