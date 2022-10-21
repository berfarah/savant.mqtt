package savant

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

func fromJSON(filepath string) ([]*Light, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	var out struct{ Lighting []savantLight }
	if err := json.Unmarshal([]byte(bytes), &out); err != nil {
		return nil, err
	}

	lighting := make([]*Light, len(out.Lighting))
	for i := range out.Lighting {
		lighting[i] = &Light{savant: out.Lighting[i]}
	}

	return lighting, nil
}

type savantLight struct {
	// State
	// `sclibridge readstate [RPMStateName]`
	State1 struct {
		Component    string
		RPMStateName string `json:"RPMStateName"`
	}

	ButtonLabel string `json:"Button Label"`

	// Service Request Parameters
	// ServiceRequest: [zone] [controller] [logicalComponent] [variant = 1] [serviceId] [command]
	// "Address1" [Address1] "Address2" [Address2] [StateChange] [Value]
	ControllerZone   string `json:"Controller Zone"`
	Controller       string
	LogicalComponent string `json:"Logical Component"`
	ServiceID        string `json:"ServiceID"`
	Command          string
	Address1         string
	Address2         string
}

func (s savantLight) buildArgs(command string) []string {
	return []string{
		s.ControllerZone,
		s.Controller,
		s.LogicalComponent,
		"1", // variant
		s.ServiceID,
		command,
		"Address1",
		s.Address1,
		"Address2",
		s.Address2,
	}
}

func (s savantLight) isDimmer() bool {
	return s.Command == "DimmerSet"
}

func (s savantLight) onArgs() []string {
	args := s.buildArgs("SwitchOn")
	if s.isDimmer() {
		args = append(s.buildArgs("DimmerSet"), "DimmerLevel", "100")
	}

	return args
}

func (s savantLight) offArgs() []string {
	args := s.buildArgs("SwitchOff")
	if s.isDimmer() {
		args = append(s.buildArgs("DimmerSet"), "DimmerLevel", "0")
	}

	return args
}

func (s savantLight) setArgs(level int) []string {
	return append(s.buildArgs("DimmerSet"), "DimmerLevel", strconv.Itoa(level))
}

type Light struct {
	savant savantLight
	// Level is the brighness value from 0-100 (only ends for switches)
	Level int
}

// Name is the user-facing name for the light (eg: Living Lights)
func (l Light) Name() string {
	return l.savant.ControllerZone + " " + l.savant.ButtonLabel
}

// ID is the machine-facing name for the light (eg: 001_01)
func (l Light) ID() string {
	return l.savant.Address1 + "_" + l.savant.Address2
}

func (l Light) IsDimmer() bool {
	return l.savant.isDimmer()
}

func (l Light) stateName() string {
	return l.savant.State1.RPMStateName
}

// State returns the light on/off state
func (l Light) State() string {
	if l.Level > 0 {
		return "ON"
	}
	return "OFF"
}
