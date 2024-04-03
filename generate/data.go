package main

import "encoding/json"

const DataFilename = "dataTableInfo.plist"
const StateVariablePrefix = "SetDimmerLevel_1_"

type Data struct {
	Lighting struct {
		// ignore Zone Map, Savant App Zone Map
		Lighting []LightingItem `plist:"Lighting"`
	} `plist:"Lighting"`
}

type LightingItem struct {
	Address1 string `plist:"Address1" json:"address_1"` // eg: 001, hex code

	ServiceID        string `plist:"ServiceID" json:"service_id"`                // eg: SVC_ENV_LIGHTING
	Controller       string `plist:"Controller" json:"controller"`               // eg: Farah SEA
	LogicalComponent string `plist:"Logical Component" json:"logical_component"` // eg: RacePointMedia_host
	Zone             string `plist:"Controller Zone" json:"zone"`                // eg: Dining
	Name             string `plist:"Button Label" json:"button_name"`            // eg: Lights, Pendant
	Entity           string `plist:"Entity" json:"entity"`                       // eg: Dimmer, Switch
	State1           struct {
		Name string `plist:"RPMStateName" json:"state_name"` // eg: Farah SEA.RacePointMedia_host.CurrentDimmerLevel_1_001
	}
}

func (l LightingItem) WriteStateName() string {
	return StateVariablePrefix + l.Address1
}

func (l LightingItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID             string `json:"id"`
		Zone           string `json:"zone"`
		Name           string `json:"name"`
		IsDimmer       bool   `json:"is_dimmer"`
		ReadStateName  string `json:"read_state_name"`
		WriteStateName string `json:"write_state_name"`
	}{
		ID:             l.Address1 + "_1",
		Zone:           l.Zone,
		Name:           l.Zone + " " + l.Name,
		IsDimmer:       l.Entity == "Dimmer",
		ReadStateName:  l.State1.Name,
		WriteStateName: "userDefined." + l.WriteStateName(),
	})
}
