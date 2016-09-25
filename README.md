# Client-side RFID authentication (and general access control) for Pumping Station: One

This system is built on a BeagleBone Black and Sparkfun's USB board for RFID
readers for hardware, and PS1Auth's server-side RFID authentication on the
backend.

Client-side software is written in Go. 


Making things work
-----
**The BeagleBone**

This software is targeted at Debian Jessie. From the default image, update `/etc/apt/sources.list` to point at 'jessie' instead of 'wheezy'

`sudo apt-get update`  
`sudo apt-get dist-upgrade`

This while take a while.

**Required Packages**

*Debian*  
`sudo apt-get install golang`  

Set up your Go environment. See (Getting Started)[https://golang.org/doc/install]
for details on getting started with Go.

*Golang*

Cloning this repo to your BBB (and building an executable) is the next step.

Run `go get github.com/loansindi/ps1rfid`. This should pull in all of the
dependencies for you.

`go get github.com/hybridgroup/gobot`  
`go get github.com/hybridgroup/gobot/platforms/beaglebone`  
`go get github.com/hybridgroup/gobot/platforms/gpio`  
`go get github.com/tarm/goserial` 
`go get github.com/BurntSushi/toml`

Navigate to the ps1rfid directory.


Build and run the code in test mode:

`go build main.go config.go beagleboneRobot.go dummyRobot.go server.go`
`./main -testMode=true`

or (if you are lazy)

`ls *.go | xargs | go build`
`./ps1rfid -testMode`

*Configuring the Server*

Service settings can be passed in a config file. An example is provided as
config.toml-sample. If the `-config` option is not used, it will use the defaults.

Default settings

  * BoltPath: "rfid-tags.db"
  * RFIDurl: "https://members.pumpingstationone.org/rfid/check/FrontDoor"
  * ToggleDuration: 5
  * TogglePin: "P9_11"
  * SerialName: "/dev/ttyUSB0"
  * SerialBaud: 9600

*Configuring the BeagleBone*

The primary challenges:

* Launching application on startup of the BBB
* Enabling non-root access to the GPIO

First, you'll want to create a user to run the service. I created a user 'rfid' 

`$ sudo adduser rfid`

Go ahead and also create a group 'gpio'

`$ sudo addgroup gpio`

Here's the service file I'm using to launch the application:

```
[Unit]
Description=Rfid entry service
After=udev.service

[Service]
Type=simple
User=rfid
Group=rfid
ExecStart=$GOPATH/src/github.com/loansindi/ps1rfid/ps1frid
WorkingDirectory=/srv/rfid
Restart=always
RestartSec=5
[Install]
WantedBy=multi-user.target
```

Drop this (with appropriate paths) in `/etc/systemd/system/rfid.service`

The udev rules go in `/etc/udev/rules.d/90-gpio.rules` :

```
KERNEL=="gpio*", SUBSYSTEM=="gpio", ACTION=="add", PROGRAM="/bin/sh -c 'chown -R rfid:gpio /sys/class/gpio'"
KERNEL=="gpio*", SUBSYSTEM=="gpio", ACTION=="add", PROGRAM="/bin/sh -c 'chown -R rfid:gpio /sys/devices/virtual/gpio/'"
KERNEL=="gpio*", SUBSYSTEM=="gpio", ACTION=="add", PROGRAM="/bin/sh -c 'chown -R rfid:gpio /sys/devices/bone_capemgr.9'"
```

This should allow non-elevated access to the GPIO, which is used both for toggling a strike plate and for an internet doorbell.

