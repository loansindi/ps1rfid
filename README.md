# Client-side RFID authentication (and general access control) for Pumping Station: One

This system is built on a BeagleBone Black and Sparkfun's USB board for RFID readers for hardware, and PS1Auth's server-side RFID authentication on the backend.

Client-side software is written in Go. 


Making things work
-----
Cloning this repo to your BBB is the first step. The primary challenges:

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
group=rfid
ExecStart=/home/derek/go/src/github.com/loansindi/ps1rfid/ps1rfid
WorkingDirectory=/srv/rfid
Restart=always

[Install]
WantedBy=multi-user.target
```

Drop this in `/etc/systemd/system/rfid.service`

The udev rules go in `/etc/udev/rules.d/90-gpio.rules` :

```
KERNEL=="gpio*", SUBSYSTEM=="gpio", ACTION=="add", PROGRAM="/bin/sh -c 'chown -R rfid:gpio /sys/class/gpio'"
KERNEL=="gpio*", SUBSYSTEM=="gpio", ACTION=="add", PROGRAM="/bin/sh -c 'chown -R rfid:gpio /sys/devices/virtual/gpio/'"
```

This should allow non-elevated ccess to the GPIO.


=======
