# fxpixel [![Build and release](https://github.com/danesparza/fxpixel/actions/workflows/release.yaml/badge.svg)](https://github.com/danesparza/fxpixel/actions/workflows/release.yaml)
REST service for RGB(W) LED lighting effects on demand from Raspberry Pi. Made with ❤️ for makers, DIY craftsmen, and professional prop designers everywhere.

## Installation
### Prerequisites
Install Raspberry Pi OS.  For best results, use the [Raspberry Pi imager](https://www.raspberrypi.com/software/)
and pick 'Raspberry Pi OS (other)' and then 'Raspberry Pi OS Lite (64-bit)'.

Install the package repo (you only need to do this once per machine)
```
wget https://packages.cagedtornado.com/prereq.sh -O - | sh
```

### Package installation
Install the package
```
sudo apt install fxpixel
```

### Timeline color continuity

The LED strip is initialized once and reused across timelines. A fade at the
start of a new timeline begins from the previous timeline's final pixel colors.
Explicit stop requests still clear the LEDs. Restart fxpixel after changing
hardware settings (GPIO, LED count, pixel order, or number of color channels).
Concurrent timelines still share the same physical output; they are not queued.
