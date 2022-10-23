# Savant MQTT Server

This server is designed to run on Savant's smart hosts and acts as an MQTT
client and state polling service.

It's set up to work with [Homeassistant's MQTT discovery](https://www.home-assistant.io/docs/mqtt/discovery/), where it configures devices as [Light Entities](https://www.home-assistant.io/integrations/light.mqtt).

## Setup

### Necessary files for the service

1. Get your lights registry from Savant by exporting loads (plist) and
   converting them with `./transform [load.plist] > [oldfile]`
2. Set up your MQTT broker in your environment
3. Create a `savantmqtt.conf` file (see [config](./config/config.go) for
   details on settings.

### Configuring Savant

Because managing state is faster than making service requests, we'll want to use
triggers to update our lights.

#### Required Workflows for Triggers

Steps:
- View Services
- Pick the service where you need lights to match state
- Create a new Service Request (naming convention: `DimmerSetVariable_{Address2}_{Address1}`
- Double click the service to open in automator
- Savant Action Argument Setter as step 1 with DimmerLevel from State Center and
    Value of `userDefined.SetDimmerLevel_{Address2}_{Address1}`
- Main action with DimmerSet on Lighting Controller Source

#### Triggers

You can generate the triggers with: `./converttrigger/transform > tmp.json`

```bin/bash
./converttrigger/transform [load.plist] > tmp.json
go run ./converttrigger/main.go tmp.json
# Import tmp.plist as into triggers - you have to do this in savant
# Clean up
rm tmp.json tmp.plist
```

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
