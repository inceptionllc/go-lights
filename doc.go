// Copyright 2015 Inception LLC.

/*

Package lights provides standardized shared functionality used across
all Go-based agent software running on light controller devices. Most
of the agents behavior has a default setting that will be "standard"
when running on devices. However, all the settings can be modified
by setting a variety of environment variables. The current environment
variables understood by agents include:

### INC_DEVICE_ID

Sets the device ID. The default device ID is formed using the hardware
(MAC) address of the wlan0 or en0 network interfaces if they exist.
For Intel Edison, wlan0 will be a reliable unique ID for each module.
Override the value of the device ID by setting the `INC_DEVICE_ID` variable
to any arbitrary string. Device IDs are used for creating nsq topic names
and reported in logs.

### INC_NSQD

Sets the nsqd address to use when agents post messages. The default is
to use the address `127.0.0.1:4150` assuming that nsqd is running locally
(which is typical of an nsq deployment). You can override the setting using
the INC_NSQD_TCP environment variable setting `host:port` as the value.

### INC_NSQLOOKUPD

Sets the nsqlookupd address to use when agents want to subscribe to messages.
The default is to use the address `nsqlookupd.local:4161` which could be set
using either rendezvous or editing the /etc/hosts file. In some cases,
it is better to set the nsqlookupd address manually. The environment variable
can contain more than one address separated by commas. For example:

	INC_NSQLOOKUPD=192.168.0.64:4161,192.168.0.61:4161
*/
package lights
