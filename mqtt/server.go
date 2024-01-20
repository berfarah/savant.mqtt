package mqtt

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/berfarah/savant.mqtt/config"
	"github.com/berfarah/savant.mqtt/savant"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Server struct {
	connected     chan bool
	connectedOnce int
	config        *config.Config
	Client        mqtt.Client
	Manager       savant.LightsManager
}

func New(config *config.Config, manager savant.LightsManager) *Server {
	return &Server{connected: make(chan bool), config: config, Manager: manager}
}

func (s Server) buildTopic(light *savant.Light, parts ...string) string {
	parts = append([]string{s.config.TopicPrefix, "light", s.config.TopicNodeID, light.ID}, parts...)
	return strings.Join(parts, "/")
}

func (s Server) topicToID(topic string) (id string) {
	id = strings.TrimPrefix(topic, strings.Join([]string{s.config.TopicPrefix, "light", s.config.TopicNodeID, ""}, "/"))
	return strings.TrimSuffix(id, "/set")
}

func (s Server) discoverySetup() {
	tokens := make(map[string]mqtt.Token)

	for id, light := range s.Manager.Lights {
		discoveryPayload := map[string]interface{}{
			"name":          light.Name,
			"schema":        "json",
			"state_topic":   s.buildTopic(light),
			"command_topic": s.buildTopic(light, "set"),
			"unique_id":     light.ID,
			"device": map[string]interface{}{
				"name":           light.Name,
				"identifiers":    light.ID,
				"manufacturer":   "Savant",
				"model":          "Light",
				"sw_version":     "savant.mqtt",
				"suggested_area": light.Zone,
			},
		}

		if light.IsDimmer {
			discoveryPayload["brightness"] = true
			discoveryPayload["color_mode"] = true
			discoveryPayload["brightness_scale"] = 100
			discoveryPayload["supported_color_modes"] = []string{"brightness"}
		}

		payload, err := json.Marshal(discoveryPayload)
		if err != nil {
			log.Printf("ERROR: Failed to convert into JSON %v: %v\n", id, err)
			continue
		}

		tokens[id] = s.Client.Publish(s.buildTopic(light, "config"), 0, true, payload)
	}

	for id, token := range tokens {
		if token.Wait(); token.Error() != nil {
			log.Printf("ERROR: Failed to create discovery channel for %s: %v\n", id, token.Error())
		}
	}
}

func (s Server) subscriptions() {
	tokens := make(map[string]mqtt.Token)
	for id, light := range s.Manager.Lights {
		tokens[id] = s.Client.Subscribe(s.buildTopic(light, "set"), 0, nil)
	}

	for id, token := range tokens {
		if token.Wait(); token.Error() != nil {
			log.Printf("ERROR: Failed to subscribe to %s: %v\n", id, token.Error())
		}
	}
}

func (s *Server) OnConnect(client mqtt.Client) {
	s.Client = client
	log.Println("DEBUG: Connected!")
	s.discoverySetup()
	s.subscriptions()
	if s.connectedOnce == 0 {
		close(s.connected)
	}
	s.connectedOnce += 1
}

type mqttPayload struct {
	State      string `json:"state"`
	Brightness int    `json:"brightness"`
}

func (s *Server) Run(ctx context.Context) {
	log.Println("DEBUG: Connecting...")
	<-s.connected

	log.Println("DEBUG: Starting polling cycle")
	s.Manager.Poll(ctx, func(event savant.StateChange) {
		light, ok := s.Manager.Lights[event.ID]
		if !ok {
			log.Println("Couldn't locate light with ID", event.ID)
			return
		}

		payload := mqttPayload{State: event.State(), Brightness: event.Level}
		b, err := json.Marshal(payload)
		if err != nil {
			log.Printf("ERROR: Failed to convert into JSON %v: %v\n", event.ID, err)
			return
		}

		token := s.Client.Publish(s.buildTopic(light), 0, true, b)
		go func() {
			if token.Wait(); token.Error() != nil {
				log.Printf("ERROR: Failed to publish message to %v: %v\n", event.ID, token.Error())
			}
		}()
	})

	log.Println("DEBUG: Stopping polling")
}

func (s Server) Handler(client mqtt.Client, msg mqtt.Message) {
	id := s.topicToID(msg.Topic())
	light, ok := s.Manager.Lights[id]
	if !ok {
		log.Println("Couldn't locate light with ID", id)
		return
	}

	var payload mqttPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("Error: Failed to decode payload for %v: %v\n", id, err)
		return
	}

	if light.IsDimmer {
		if payload.State == "OFF" {
			payload.Brightness = 0
		}

		if payload.State == "ON" && payload.Brightness == 0 {
			payload.Brightness = 100
		}

		s.Manager.Set(id, payload.Brightness)
		return
	}

	if payload.State == "ON" {
		s.Manager.TurnOn(id)
	}

	if payload.State == "OFF" {
		s.Manager.TurnOff(id)
	}
}
