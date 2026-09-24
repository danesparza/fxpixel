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


## Network discovery (fxcontroller)

Building and testing requires Go 1.25 or newer.

The service advertises its HTTP API through Zeroconf (mDNS/DNS-SD) by
default. A future fxcontroller can browse **`_fx._tcp` in `local.`** once
to find fxaudio, fxpixel, fxdmx, and fxtrigger. No controller is required
to run the services, and direct HTTP API access remains available if
multicast registration fails (a warning is logged).

```yaml
discovery:
  enabled: true
  name: "" # Optional friendly instance name, 1–63 UTF-8 bytes
  id: ""   # Optional unique, stable installation ID
```

The default name is `service-hostname-port`; the default ID is
`service:hostname:port`. Give each installation a distinct `discovery.id`
if identity must survive hostname or port changes. Names must also be
unique on the LAN. TXT values must fit the DNS-SD 255-byte record limit.

The discovery contract is the same across all four projects:

| DNS record | Meaning |
| --- | --- |
| SRV / A / AAAA | HTTP hostname, actual listening port, and addresses |
| TXT `txtvers=1` | Discovery metadata format version |
| TXT `service` | `fxaudio`, `fxpixel`, `fxdmx`, or `fxtrigger` |
| TXT `id` | Installation identity |
| TXT `api=v1` | API version |
| TXT `scheme=http` | API transport |
| TXT `path=/v1` | API base path |

Controllers should use the resolved SRV port and address, dispatch by
`service`, check supported metadata/API versions, and verify API reachability
separately. Discovery is not authentication; treat advertisements as untrusted
network input. This adds service advertisements, not a controller or UI.

Advertisements start only after the HTTP listener binds and are withdrawn
on SIGINT/SIGTERM or when the HTTP server stops. Discovery uses UDP 5353
multicast on the local network segment; firewalls, Wi-Fi client isolation,
VLANs, and container networking can prevent discovery. Containers typically
need host networking (where supported) or an mDNS relay. No internet service
or separately installed Bonjour/Avahi daemon is required.

On macOS, inspect live advertisements with:

```sh
dns-sd -B _fx._tcp local.
# Resolve one of the instance names returned above:
dns-sd -L "INSTANCE NAME" _fx._tcp local.
```

On Linux with Avahi tools: `avahi-browse -rt _fx._tcp`.
