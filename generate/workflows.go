package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const savantPath = "Library/Application Support/Savant"

func getSavantDir(sys SystemInfo, doc DocumentInfo) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	savantSpace := path.Join(homeDir, savantPath)

	entries, err := os.ReadDir(savantSpace)
	if err != nil {
		return "", err
	}

	savantVersions := []string{}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".da Vinci "+doc.VersionKey+" "+doc.BuildNumber) {
			savantVersions = []string{e.Name()}
			break
		}
		if strings.HasPrefix(e.Name(), ".da Vinci") {
			savantVersions = append(savantVersions, e.Name())
		}
	}

	if len(savantVersions) == 0 {
		return "", fmt.Errorf("No version of Savant found in %s.\nMake sure you've installed one in Savant Application Manager.", savantSpace)
	}

	version := savantVersions[0]
	selection := 0

	if len(savantVersions) > 1 {
		if err := retrySelection(3, func() (err error) {
			selection, err = selectItem("Version", savantVersions)
			return err
		}); err != nil {
			return "", err
		}
		version = savantVersions[selection]
	}

	return path.Join(savantSpace, version, "RPMInstallLink/Library"), nil

}

func getSavantCacheDir(savantSpace string) (string, error) {
	cachesDir := path.Join(savantSpace, "Application Support/RacePointMedia/Caches/com.savantav.RacePointBlueprint")

	entries, err := os.ReadDir(cachesDir)
	if err != nil {
		return "", err
	}

	caches := []string{}
	for _, e := range entries {
		if e.Name() == ".DS_Store" {
			continue
		}
		caches = append(caches, e.Name())
	}

	if len(caches) == 0 {
		return "", fmt.Errorf("No caches of active Savant sessions found in %s.\nMake sure you have Blueprint running with your active project.", cachesDir)
	}

	cache := caches[0]
	cacheSelection := 0

	if len(caches) > 1 {
		if err := retrySelection(3, func() (err error) {
			cacheSelection, err = selectItem("Project Cache", caches)
			return err
		}); err != nil {
			return "", err
		}
		cache = caches[cacheSelection]
	}

	return path.Join(cachesDir, cache), nil
}

func retrySelection(attempts int, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if err = f(); err == nil {
			return nil
		}
	}
	return err
}

func selectItem(text string, items []string) (int, error) {
	fmt.Println("Select", text+":")
	for i, v := range items {
		fmt.Println(strconv.Itoa(i+1), ") ", v)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("An error occurred while reading input. Please try again.")
	}

	input = strings.TrimSpace(input)
	selection, err := strconv.Atoi(input)

	if err != nil {
		return 0, fmt.Errorf("Please input a numbered selection")
	}

	if selection > len(items) {
		return 0, fmt.Errorf("Please select a valid version")
	}

	return selection - 1, nil
}

func createWorkflow(dir, caches string, l LightingItem) error {
	workflowDir := path.Join(caches, "workflowsTemporaryFolder",
		l.Zone,
		l.Controller+"-"+l.LogicalComponent+"-1-"+l.ServiceID,
		WorkflowPrefix+l.Address1+".workflow",
		"Contents",
	)

	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return err
	}

	workflowText := generateWorkflowText(dir, l)
	if err := os.WriteFile(path.Join(workflowDir, "document.wflow"), []byte(workflowText), 0644); err != nil {
		return err
	}

	infoText := generateInfoText(l.WriteStateName())
	if err := os.WriteFile(path.Join(workflowDir, "Info.plist"), []byte(infoText), 0644); err != nil {
		return err
	}

	return nil
}

func generateInfoText(name string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>` + name + `</string>
</dict>
</plist>
`
}

func generateWorkflowText(dir string, l LightingItem) string {
	argUUID := strings.ToUpper(uuid.NewString())
	argInputUUID := strings.ToUpper(uuid.NewString())
	argOutputUUID := strings.ToUpper(uuid.NewString())

	// hard coded for now
	saPath := "/Users/bernardofarah/Library/Automator/Farah SEA.action"
	saUUID := strings.ToUpper(uuid.NewString())
	saInputUUID := strings.ToUpper(uuid.NewString())
	saOutputUUID := strings.ToUpper(uuid.NewString())

	connectorUUID := strings.ToUpper(uuid.NewString())

	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AMApplicationBuild</key>
	<string>512</string>
	<key>AMApplicationVersion</key>
	<string>2.10</string>
	<key>AMDocumentVersion</key>
	<string>2</string>
	<key>actions</key>
	<array>

		<dict>
			<key>action</key>
			<dict>
				<key>AMAccepts</key>
				<dict>
					<key>Container</key>
					<string>List</string>
					<key>Optional</key>
					<true/>
					<key>Types</key>
					<array>
						<string>com.apple.cocoa.path</string>
					</array>
				</dict>
				<key>AMActionVersion</key>
				<string>722</string>
				<key>AMApplication</key>
				<array>
					<string>Unknown</string>
				</array>
				<key>AMParameterProperties</key>
				<dict>
					<key>RPMClassName</key>
					<dict/>
					<key>argName</key>
					<dict/>
					<key>argValue</key>
					<dict/>
					<key>valueSource</key>
					<dict/>
				</dict>
				<key>AMProvides</key>
				<dict>
					<key>Container</key>
					<string>List</string>
					<key>Types</key>
					<array>
						<string>com.apple.cocoa.path</string>
					</array>
				</dict>
				<key>AMRequiredResources</key>
				<array/>
				<key>ActionBundlePath</key>
				<string>` + dir + `/Automator/Savant Action Argument Set.action</string>
				<key>ActionName</key>
				<string>Savant Action Argument Setter</string>
				<key>ActionParameters</key>
				<dict>
					<key>RPMClassName</key>
					<string>rpmActionArgSet</string>
					<key>argName</key>
					<string>DimmerLevel</string>
					<key>argValue</key>
					<string>userDefined.` + l.WriteStateName() + `</string>
					<key>valueSource</key>
					<string>State Center</string>
				</dict>
				<key>BundleIdentifier</key>
				<string>com.savangtav.Automator.Savant_Action_Argument_Set</string>
				<key>CFBundleVersion</key>
				<string>722</string>
				<key>CanShowSelectedItemsWhenRun</key>
				<true/>
				<key>CanShowWhenRun</key>
				<true/>
				<key>Category</key>
				<array>
					<string>Savant</string>
				</array>
				<key>Class Name</key>
				<string>AMRemoteAction</string>
				<key>InputUUID</key>
				<string>` + argInputUUID + `</string>
				<key>Keywords</key>
				<array/>
				<key>OutputUUID</key>
				<string>` + argOutputUUID + `</string>
				<key>UUID</key>
				<string>` + argUUID + `</string>
				<key>UnlocalizedApplications</key>
				<array>
					<string>Savant</string>
				</array>
				<key>arguments</key>
				<dict/>
				<key>isViewVisible</key>
				<integer>1</integer>
				<key>location</key>
				<string>408.500000:197.000000</string>
				<key>nibPath</key>
				<string>` + dir + `/Automator/Savant Action Argument Set.action/Contents/Resources/English.lproj/main.nib</string>
			</dict>
			<key>isViewVisible</key>
			<integer>1</integer>
		</dict>

		<dict>
			<key>action</key>
			<dict>
				<key>AMAccepts</key>
				<dict>
					<key>Container</key>
					<string>List</string>
					<key>Optional</key>
					<false/>
					<key>Types</key>
					<array/>
				</dict>
				<key>AMActionVersion</key>
				<string>722</string>
				<key>AMApplication</key>
				<array>
					<string>Unknown</string>
				</array>
				<key>AMParameterProperties</key>
				<dict>
					<key>RPMClassName</key>
					<dict/>
					<key>actionName</key>
					<dict/>
					<key>arguments</key>
					<dict/>
					<key>componentID</key>
					<dict/>
					<key>logicalComponentName</key>
					<dict/>
					<key>requestingServiceName</key>
					<dict/>
					<key>requestingServiceVariant</key>
					<dict/>
					<key>resourceName</key>
					<dict/>
					<key>workflowBlocking</key>
					<dict/>
				</dict>
				<key>AMProvides</key>
				<dict>
					<key>Container</key>
					<string>List</string>
					<key>Types</key>
					<array/>
				</dict>
				<key>AMRequiredResources</key>
				<array/>
				<key>ActionBundlePath</key>
				<string>` + saPath + `</string>
				<key>ActionName</key>
				<string>` + l.Controller + `</string>
				<key>ActionParameters</key>
				<dict>
					<key>RPMClassName</key>
					<string>avAction</string>
					<key>actionName</key>
					<string>DimmerSet</string>
					<key>arguments</key>
					<array>
						<dict>
							<key>argDescription</key>
							<string>Module Address that contains Load (000 - 3FF)</string>
							<key>argName</key>
							<string>Address1</string>
							<key>argValue</key>
							<string>` + l.Address1 + `</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Load Number on Module Addressed in Address1 (0-7)</string>
							<key>argName</key>
							<string>Address2</string>
							<key>argValue</key>
							<string>1</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>not used</string>
							<key>argName</key>
							<string>Address3</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>not used</string>
							<key>argName</key>
							<string>Address4</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>not used</string>
							<key>argName</key>
							<string>Address5</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>not used</string>
							<key>argName</key>
							<string>Address6</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the Desired Value (Infinite Color,Tunable White,Warm Glow,etc...)</string>
							<key>argName</key>
							<string>Curve</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>not used</string>
							<key>argName</key>
							<string>DelayTime</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the Desired Value (0 - 100), %</string>
							<key>argName</key>
							<string>DimmerLevel</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the duration (in seconds) for the transtion to the target dimmer level. (0 - 99)</string>
							<key>argName</key>
							<string>FadeTime</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the Desired Value (0-255)</string>
							<key>argName</key>
							<string>bleColorBlue</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the Desired Value (0-255)</string>
							<key>argName</key>
							<string>bleColorGreen</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the Desired Value (0-255)</string>
							<key>argName</key>
							<string>bleColorRed</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the Desired Value (0-255)</string>
							<key>argName</key>
							<string>bleColorWhite</string>
						</dict>
						<dict>
							<key>argDescription</key>
							<string>Enter the Desired Value (2200-6000)</string>
							<key>argName</key>
							<string>kelvin</string>
						</dict>
					</array>
					<key>componentID</key>
					<string>` + l.Controller + `</string>
					<key>logicalComponentName</key>
					<string>RacePointMedia_host</string>
					<key>requestingServiceName</key>
					<string>SVC_ENV_LIGHTING</string>
					<key>requestingServiceVariant</key>
					<string>1</string>
					<key>resourceName</key>
					<string>Lighting Controller Source</string>
					<key>workflowBlocking</key>
					<false/>
				</dict>
				<key>BundleIdentifier</key>
				<string>com.racepointmedia.Automator.` + l.Controller + `</string>
				<key>CFBundleVersion</key>
				<string>722</string>
				<key>CanShowSelectedItemsWhenRun</key>
				<false/>
				<key>CanShowWhenRun</key>
				<true/>
				<key>Category</key>
				<array>
					<string>Savant</string>
				</array>
				<key>Class Name</key>
				<string>AMRemoteAction</string>
				<key>InputUUID</key>
				<string>` + saInputUUID + `</string>
				<key>Keywords</key>
				<array>
					<string>Component Control</string>
					<string>Savant</string>
				</array>
				<key>OutputUUID</key>
				<string>` + saOutputUUID + `</string>
				<key>UUID</key>
				<string>` + saUUID + `</string>
				<key>UnlocalizedApplications</key>
				<array>
					<string>Savant</string>
				</array>
				<key>arguments</key>
				<dict/>
				<key>isViewVisible</key>
				<integer>1</integer>
				<key>location</key>
				<string>408.500000:478.000000</string>
				<key>nibPath</key>
				<string>` + saPath + `/Contents/Resources/English.lproj/main.nib</string>
			</dict>
			<key>isViewVisible</key>
			<integer>1</integer>
		</dict>

	</array>

	<key>connectors</key>
	<dict>
		<key>` + connectorUUID + `</key>
		<dict>
			<key>from</key>
			<string>` + argUUID + ` - ` + argUUID + `</string>
			<key>to</key>
			<string>` + saUUID + ` - ` + saUUID + `</string>
		</dict>
	</dict>

	<key>workflowMetaData</key>
	<dict>
		<key>workflowTypeIdentifier</key>
		<string>com.apple.Automator.workflow</string>
	</dict>
</dict>
</plist>
`
}
