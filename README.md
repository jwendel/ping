# ping

**THIS IS A WORK IN PROGRESS.  I RECOMMEND AGAINST USING IT AT THIS TIME.**

The ping library for Golang (Go) provides ICMP ECHO functionality.

### Features

* Supports multiple target hosts.
* Results come back through a channel.


## Requirements
Ping needs a raw socket to work, which normally means it needs to be run as root.  There are various ways to get this working on different platforms.

#### Linux
On Debian/Ubuntu can install the **libcap2-bin** package, which adds the command **setcap**.  Use setcap to set *CAP\_NET\_RAW* capability on your binary, ping will be able to create it's socket without root permission.  Example:

* setcap cap\_net\_raw=ep ./ping

#### Mac OSX or Linux
By default when you call Start() on a Pinger, it will attempt to call **setuid(0)** to make itself root before grabbing the raw socket.  It will set itself back to the running users UID after it grabs the socket.

Example:

* sudo chown root ping
* sudo chmod u+s ping
* ./ping golang.org

#### Windows
This requires administrator privilages to use, I'm not aware of any easy workarounds.  There are [some windows APIs](http://msdn.microsoft.com/en-us/library/windows/desktop/aa366050.aspx) that allow ping without admin, but it is not something I've looked using for this library.  So run your app as administrator on Windows.

## Known Issues

* While multiple instances of Pinger can be created, creating 2 instances than pinging the same IP address from both will have odd behavior.
* Changing configuraiton settings while Pinger is running is untested.  Call Stop before changing settings.

## TODO:

* add Stop command.

