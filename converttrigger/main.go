package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const start = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>RPMExportDataTypeKey</key>
	<string>RPMTriggersListType</string>
	<key>RPMExportListKey</key>
	<array>
		<dict>
			<key>Group Name</key>
			<string>Lights</string>
			<key>Table</key>
			<array>
`

const end = `
			</array>
		</dict>
	</array>
	<key>RPMExportVersionKey</key>
	<integer>1</integer>
</dict>
</plist>
`

type Input struct {
	Zone       string `json:"zone"`
	ButtonName string `json:"button_name"`
	Address1   string `json:"address_1"`
	IsDimmer   bool   `json:"is_dimmer"`
}

func dimmerText(i Input) string {
	return `
				<dict>
					<key>Else Actions</key>
					<array/>
					<key>Name</key>
					<string>` + i.Zone + " " + i.ButtonName + `</string>
					<key>State Match</key>
					<array/>
					<key>Then Actions</key>
					<array>
						<dict>
							<key>Service Request</key>
							<dict>
								<key>RPMServicePasteboardServiceStateKey</key>
								<string></string>
								<key>Request</key>
								<string>DimmerSetVariable_1_` + i.Address1 + `</string>
								<key>Request Data</key>
								<dict/>
								<key>Service Type</key>
								<string>SVC_ENV_LIGHTING</string>
								<key>Source Component Name</key>
								<string>Farah SEA</string>
								<key>Source Logical Component Name</key>
								<string>RacePointMedia_host</string>
								<key>Variant ID</key>
								<string>1</string>
								<key>Zone</key>
								<string>` + i.Zone + `</string>
							</dict>
						</dict>
					</array>
					<key>Trigger States</key>
					<array>
						<dict>
							<key>RPMStateName</key>
							<string>userDefined.SetDimmerLevel_1_` + i.Address1 + `</string>
							<key>RPMStateType</key>
							<string>RPMBasicStateName</string>
							<key>State Name</key>
							<string>SetDimmerLevel_1_` + i.Address1 + `</string>
							<key>State Scope</key>
							<string>userDefined</string>
							<key>State Type</key>
							<string>integer</string>
							<key>identifiers</key>
							<array/>
							<key>scope</key>
							<string>userDefined</string>
							<key>stateName</key>
							<string>SetDimmerLevel_1_` + i.Address1 + `</string>
						</dict>
					</array>
				</dict>
	`
}

func switchText(i Input) string {
	return `
				<dict>
					<key>Else Actions</key>
					<array>
						<dict>
							<key>Service Request</key>
							<dict>
								<key>RPMServicePasteboardServiceStateKey</key>
								<string>inactive</string>
								<key>Request</key>
								<string>SwitchOn</string>
								<key>Request Data</key>
								<dict>
									<key>Address1</key>
									<string>` + i.Address1 + `</string>
									<key>Address2</key>
									<string>1</string>
								</dict>
								<key>Service Type</key>
								<string>SVC_ENV_LIGHTING</string>
								<key>Source Component Name</key>
								<string>Farah SEA</string>
								<key>Source Logical Component Name</key>
								<string>RacePointMedia_host</string>
								<key>Variant ID</key>
								<string>1</string>
								<key>Zone</key>
								<string>` + i.Zone + `</string>
							</dict>
						</dict>
					</array>
					<key>Name</key>
					<string>` + i.Zone + " " + i.ButtonName + `</string>
					<key>State Match</key>
					<array>
						<dict>
							<key>Match Data</key>
							<string>0</string>
							<key>Match Logic</key>
							<string>Equal</string>
							<key>RPMStateName</key>
							<string>userDefined.SetDimmerLevel_1_` + i.Address1 + `</string>
							<key>RPMStateType</key>
							<string>RPMBasicStateName</string>
							<key>State Name</key>
							<string>SetDimmerLevel_1_` + i.Address1 + `</string>
							<key>State Scope</key>
							<string>userDefined</string>
							<key>State Type</key>
							<string>integer</string>
							<key>identifiers</key>
							<array/>
							<key>scope</key>
							<string>userDefined</string>
							<key>stateName</key>
							<string>SetDimmerLevel_1_` + i.Address1 + `</string>
						</dict>
					</array>
					<key>Then Actions</key>
					<array>
						<dict>
							<key>Service Request</key>
							<dict>
								<key>RPMServicePasteboardServiceStateKey</key>
								<string>inactive</string>
								<key>Request</key>
								<string>SwitchOff</string>
								<key>Request Data</key>
								<dict>
									<key>Address1</key>
									<string>` + i.Address1 + `</string>
									<key>Address2</key>
									<string>1</string>
								</dict>
								<key>Service Type</key>
								<string>SVC_ENV_LIGHTING</string>
								<key>Source Component Name</key>
								<string>Farah SEA</string>
								<key>Source Logical Component Name</key>
								<string>RacePointMedia_host</string>
								<key>Variant ID</key>
								<string>1</string>
								<key>Zone</key>
								<string>` + i.Zone + `</string>
							</dict>
						</dict>
					</array>
					<key>Trigger States</key>
					<array>
						<dict>
							<key>RPMStateName</key>
							<string>userDefined.SetDimmerLevel_1_` + i.Address1 + `</string>
							<key>RPMStateType</key>
							<string>RPMBasicStateName</string>
							<key>State Name</key>
							<string>SetDimmerLevel_1_` + i.Address1 + `</string>
							<key>State Scope</key>
							<string>userDefined</string>
							<key>State Type</key>
							<string>integer</string>
							<key>identifiers</key>
							<array/>
							<key>scope</key>
							<string>userDefined</string>
							<key>stateName</key>
							<string>SetDimmerLevel_1_` + i.Address1 + `</string>
						</dict>
					</array>
				</dict>
	`
}

func fromJSON(filepath string) ([]Input, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	var out []Input
	if err := json.Unmarshal([]byte(bytes), &out); err != nil {
		return nil, err
	}

	return out, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing required argument (json file)")
		os.Exit(1)
	}

	inputs, err := fromJSON(os.Args[1])
	if err != nil {
		fmt.Println("Error decoding JSON file", err.Error())
		os.Exit(1)
	}

	var str string
	str = start
	for _, input := range inputs {
		if input.IsDimmer {
			str = str + dimmerText(input)
		} else {
			str = str + switchText(input)
		}
	}
	str = str + end

	outputFile := strings.Split(os.Args[1], ".")[0] + ".plist"
	if err := os.WriteFile(outputFile, []byte(str), 0644); err != nil {
		fmt.Println("Error writing file", err.Error())
		os.Exit(1)
	}
}
