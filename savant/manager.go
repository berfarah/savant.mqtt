package savant

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/berfarah/savant.mqtt/config"
)

type LightsManagerOptions struct {
	RegistryFile    string
	PollingInterval time.Duration
}

func New(config *config.Config) (LightsManager, error) {
	lights, err := fromJSON(config.RegistryFilePath)
	if err != nil {
		return LightsManager{}, err
	}

	ids := make([]string, len(lights))
	lightsMap := make(map[string]*Light)
	for i, light := range lights {
		ids[i] = light.ID()
		lightsMap[light.ID()] = light
	}

	return LightsManager{
		config:  config,
		ids:     ids,
		Lights:  lightsMap,
		stateCh: make(chan StateChange, len(lights)*2),
	}, nil
}

type LightsManager struct {
	config      *config.Config
	ids         []string
	Lights      map[string]*Light
	lastUpdated time.Time
	stateCh     chan StateChange
}

type StateChange struct {
	ID    string
	Level int
}

// State returns the light on/off state
func (sc StateChange) State() string {
	if sc.Level > 0 {
		return "ON"
	}
	return "OFF"
}

func (l LightsManager) refreshState() error {
	stateNames := make([]string, 0, len(l.Lights))
	for _, id := range l.ids {
		stateNames = append(stateNames, l.Lights[id].stateName())
	}

	states, err := scliClient.Run("readstate", stateNames...)
	if err != nil {
		return err
	}

	for i, state := range states {
		id := l.ids[i]
		level, err := strconv.ParseFloat(state, 2)
		if err != nil {
			fmt.Println("Invalid state:", err)
		}
		l.setState(id, int(level))
	}

	return nil
}

func (l LightsManager) runPoller(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			if err := l.refreshState(); err != nil {
				fmt.Println(err.Error())
			}
		case <-ctx.Done():
			close(l.stateCh)
			return
		}
	}
}

// Poll refreshes state on a time interval by querying sclibridge in batches
func (l LightsManager) Poll(ctx context.Context, cb func(StateChange)) {
	go l.runPoller(ctx, time.Duration(l.config.PollSeconds)*time.Second)
	for state := range l.stateCh {
		l.Lights[state.ID].Level = state.Level
		cb(state)
	}
}

func (l LightsManager) setState(id string, level int) {
	l.stateCh <- StateChange{ID: id, Level: level}
}

// Set sets a custom level, turns a dimmer on
func (l LightsManager) Set(id string, level int) error {
	// Cap level at 100
	if level > 100 {
		level = 100
	}

	if _, err := scliClient.Run("servicerequest", l.Lights[id].savant.setArgs(level)...); err != nil {
		return err
	}
	l.setState(id, level)

	return nil
}

// Turn On sets Level to 100, turns the light on
func (l LightsManager) TurnOn(id string) error {
	if _, err := scliClient.Run("servicerequest", l.Lights[id].savant.onArgs()...); err != nil {
		return err
	}
	l.setState(id, 100)

	return nil
}

// Turn Off sets Level to 0, turns the light off
func (l LightsManager) TurnOff(id string) error {
	if _, err := scliClient.Run("servicerequest", l.Lights[id].savant.offArgs()...); err != nil {
		return err
	}
	l.setState(id, 0)

	return nil
}
