# Hue I'm Home

> A docker image to automatically turn on/off your Hue lights when you enter/leave your home.

It works by port-scanning your network for known open Android/iOS ports,
and changes your hue lights when your device's ports open and close on the network.

 - GitHub: https://github.com/se1exin/Hue-Im-Home
 - Docker Hub: https://hub.docker.com/r/selexin/hue-im-home

## Usage

Example for running with **docker**

```
docker run \
 -e IP_RANGE="10.1.1.11-20" \
 -e BRIDGE_IP="10.1.1.10" \
 -e SCAN_INTERVAL=10 \
 -v </path/to/appdata/config>:/config \
 --restart unless-stopped \
 selexin/hue-im-home
```

## First Start

Please ensure the Bridge Button has been pressed before starting the container.

On first start, the docker image will create a `config.json` file in the mounted `/config` volume,
and attempt to detect and register with your Hue Bridge.

The Bridge API key is persisted to the `config.json` file to avoid re-creating new 'users' on
your bridge every time the container runs.

## Parameters

Container images are configured using the following parameters passed at runtime.

| Parameter | Function |
| :----: | --- |
| `-e IP_RANGE=W.X.Y.Z` | (required) nmap style IP range - see below for explanation  |
| `-e BRIDGE_IP=W.X.Y.Z` | (optional) IP Address of your Hue Bridge - see below for explanation |
| `-e SCAN_INTERVAL=` | (optional) Time in seconds to wait before re-scanning the network. Default is 10 seconds |
| `-e DEVICE_TYPE=` | (optional) set to **"android"** or **"ios"** to limit scanning to that device type |
| `-e ON_TIME_RANGE=` | (optional) Time Range to restrict lights ON behaviour - see below for explanation |
| `-e OFF_TIME_RANGE=` | (optional) Time Range to restrict lights OFF behaviour - see below for explanation |
| `-v /config` | volume to store the required config file |

### `IP_RANGE` Environment Variable
You must tell the docker container which IP addresses it should scan for devices.

You can set this to your entire DHCP Pool (you can find this info in your router),
but keep in mind the smaller the range, the quicker the scan.

For example, on my network I know my DCHP Pool starts at `10.1.1.10`, and I am unlikely to have more
than 10 DHCP devices, so I set the `IP_RANGE` environment variable to `10.1.1.10-20` 

Other examples might include:
 - `192.168.1.1-100`
 - `192.168.1.1,2,3`
 
 For more info see https://nmap.org/book/man-target-specification.html
 
### `BRIDGE_IP` Environment Variable
The docker container will attempt to scan your network for available Hue Bridges on first start,
however if this fails to auto-discover your bridge, you can manually specify the Bridge IP Address using
the `BRIDGE_IP` environment variable.

You can find your Bridge IP in the Hue App (on your phone) under:
`Settings` > `Hue Bridges` > `(Info Icon)` > `IP address`
or look on your router for it's IP Address.

### `ON_TIME_RANGE` and `OFF_TIME_RANGE` Environment Variables
You may optionally restrict the time range in which your lights will turn ON and OFF by passing a 24 hour Time String in format `hh:mm-hh:mm`.
Example:
```
ON_TIME_RANGE=14:00-18:00
OFF_TIME_RANGE=08:00-10:00
```
The above example will only turn the lights ON between 2pm and 6pm, and only turn them back off between 8am and 10am. For example when you get home from, and leave to go to, work.

Note: Your lights will only turn on at most ONCE in this time period per-day. This is required to avoid constantly turning your lights off and on once you are already home.
E.g. if you own an iPhone, as turning the screen off and on will trigger the lights to turn off and back on. 


## Supported Architectures

This image supports multiple architectures, and utilises docker manifest for multi-platform awareness.

Simply pulling selexin/hue-im-home should retrieve the correct image for your arch, but you can also pull specific arch images via tags. 

The architectures supported by this image are:

| Architecture | Tag |
| :----: | --- |
| x86-64 | amd64-latest |
| armhf | arm32v7-latest |


## Known Issues
 - [x] iOS devices close their ports shortly after the screen locks,
which will cause the program to think that the device has left, and to turn off the lights.
 - [x] The rules for turning on/off lights are very basic at the moment. Plans are in the works to add time-based rules to avoid turning on the lights in the middle of the night, etc.

## License

MIT - see [LICENSE.md](LICENSE.md)