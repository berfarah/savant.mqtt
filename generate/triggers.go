package main

const startTriggers = `<?xml version="1.0" encoding="UTF-8"?>
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

const endTriggers = `
			</array>
		</dict>
	</array>
	<key>RPMExportVersionKey</key>
	<integer>1</integer>
</dict>
</plist>
`

const WorkflowPrefix = "DimmerSetVariable_1_"

func dimmerText(i LightingItem) string {
	return `
				<dict>
					<key>Else Actions</key>
					<array/>
					<key>Name</key>
					<string>` + i.Zone + " " + i.Name + `</string>
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
								<string>` + WorkflowPrefix + i.Address1 + `</string>
								<key>Request Data</key>
								<dict/>
								<key>Service Type</key>
								<string>` + i.ServiceID + `</string>
								<key>Source Component Name</key>
								<string>` + i.Controller + `</string>
								<key>Source Logical Component Name</key>
								<string>` + i.LogicalComponent + `</string>
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
							<string>userDefined.` + i.WriteStateName() + `</string>
							<key>RPMStateType</key>
							<string>RPMBasicStateName</string>
							<key>State Name</key>
							<string>` + i.WriteStateName() + `</string>
							<key>State Scope</key>
							<string>userDefined</string>
							<key>State Type</key>
							<string>integer</string>
							<key>identifiers</key>
							<array/>
							<key>scope</key>
							<string>userDefined</string>
							<key>stateName</key>
							<string>` + i.WriteStateName() + `</string>
						</dict>
					</array>
				</dict>
	`
}

func switchText(i LightingItem) string {
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
								<string>` + i.ServiceID + `</string>
								<key>Source Component Name</key>
								<string>` + i.Controller + `</string>
								<key>Source Logical Component Name</key>
								<string>` + i.LogicalComponent + `</string>
								<key>Variant ID</key>
								<string>1</string>
								<key>Zone</key>
								<string>` + i.Zone + `</string>
							</dict>
						</dict>
					</array>
					<key>Name</key>
					<string>` + i.Zone + " " + i.Name + `</string>
					<key>State Match</key>
					<array>
						<dict>
							<key>Match Data</key>
							<string>0</string>
							<key>Match Logic</key>
							<string>Equal</string>
							<key>RPMStateName</key>
							<string>userDefined.` + i.WriteStateName() + `</string>
							<key>RPMStateType</key>
							<string>RPMBasicStateName</string>
							<key>State Name</key>
							<string>` + i.WriteStateName() + `</string>
							<key>State Scope</key>
							<string>userDefined</string>
							<key>State Type</key>
							<string>integer</string>
							<key>identifiers</key>
							<array/>
							<key>scope</key>
							<string>userDefined</string>
							<key>stateName</key>
							<string>` + i.WriteStateName() + `</string>
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
								<string>` + i.ServiceID + `</string>
								<key>Source Component Name</key>
								<string>` + i.Controller + `</string>
								<key>Source Logical Component Name</key>
								<string>` + i.LogicalComponent + `</string>
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
							<string>userDefined.` + i.WriteStateName() + `</string>
							<key>RPMStateType</key>
							<string>RPMBasicStateName</string>
							<key>State Name</key>
							<string>` + i.WriteStateName() + `</string>
							<key>State Scope</key>
							<string>userDefined</string>
							<key>State Type</key>
							<string>integer</string>
							<key>identifiers</key>
							<array/>
							<key>scope</key>
							<string>userDefined</string>
							<key>stateName</key>
							<string>` + i.WriteStateName() + `</string>
						</dict>
					</array>
				</dict>
	`
}

func generateTriggers(lights []LightingItem) string {
	triggers := startTriggers

	for _, input := range lights {
		switch input.Entity {
		case "Dimmer":
			triggers = triggers + dimmerText(input)
		case "Switch":
			triggers = triggers + switchText(input)
		}
	}

	return triggers + endTriggers
}
