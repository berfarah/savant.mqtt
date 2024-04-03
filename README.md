# Savant MQTT Server

This server is designed to run on Savant's smart hosts and acts as an MQTT
client and state polling service.

It's set up to work with [Homeassistant's MQTT discovery](https://www.home-assistant.io/docs/mqtt/discovery/), where it configures devices as [Light Entities](https://www.home-assistant.io/integrations/light.mqtt).

## Setup

Before you get started generating files, make sure you open your blueprint
editor with the `rpmConfig` you plan on deploying. This will load the file into
cache, which is where our generator will put some needed details.

### Generate your files

    go run ./generate [path to].rpmConfig

Will create the following for you:
- `triggers.plist` for importing into Savant
- Necessary workflows for your triggers to work
- `loads.json` for referencing in your configuration file

### Steps you have to take yourself

1. Open your triggers window in the blueprint editor and import the
   `triggers.plist` in the bottom right corner
2. Create a `savantmqtt.conf` file (see [config](./config/config.go) for
   details on settings).

### Getting running

1. Run the build command
2. Copy the binary to the Savant host
3. Copy the systemctl config (lib/savant-mqtt.service) to the savant host
   under /lib/systemd/system/savant-mqtt.service
4. Run `sudo systemctl daemon-reload`
5. Enable the service via `sudo systemctl enable savant-mqtt`
6. Start the service via `sudo systemctl start savant-mqtt`

## Shout-outs

Thanks to [this guide](https://levelup.gitconnected.com/how-to-use-mqtt-with-go-89c617915774) for making it easy to get up and running with MQTT in Go!
