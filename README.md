# Client-side RFID authentication (and general access control) for Pumping Station: One

This system is built on a BeagleBone Black and Sparkfun's USB board for RFID readers for hardware, and PS1Auth's server-side RFID authentication on the backend.

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
`sudo apt-get install libzmq3`  
`sudo apt-get install golang`  

*Golang*
* github.com/pebbe/zmq4
* github.com/hybridgroup/gobot
* github.com/hybridgroup/gobot/platforms/beaglebone
* github.com/hybridgroup/gobot/platforms/gpio
* github.com/tarm/goserial

Cloning this repo to your BBB (and building an executable) is the next step. The primary challenges:

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


=======
