package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"howett.net/plist"
)

const loadsOutputFilename = "loads.json"
const triggersOutputFilename = "triggers.plist"

func main() {
	if len(os.Args) < 2 || !strings.HasSuffix(os.Args[1], "rpmConfig") {
		fmt.Println("Must provide rpmConfig as argument")
		os.Exit(1)
	}
	directory := os.Args[1]

	var sysInfo SystemInfo
	if err := decodeFile(path.Join(directory, SystemInfoFilename), &sysInfo); err != nil {
		fmt.Println("Error decoding System Info file:", err)
	}

	var docInfo DocumentInfo
	if err := decodeFile(path.Join(directory, DocumentInfoFilename), &docInfo); err != nil {
		fmt.Println("Error decoding Document Info file:", err)
	}

	var dataInfo Data
	if err := decodeFile(path.Join(directory, DataFilename), &dataInfo); err != nil {
		fmt.Println("Error decoding Data Info file:", err)
	}

	var filteredLighting []LightingItem
	for _, i := range dataInfo.Lighting.Lighting {
		if i.Entity != "Keypad Button" {
			filteredLighting = append(filteredLighting, i)
		}
	}
	dataInfo.Lighting.Lighting = filteredLighting

	fmt.Println(sysInfo, docInfo, dataInfo.Lighting.Lighting)
	workingDir, _ := os.Getwd()

	// Write loads config
	fmt.Println("Writing Loads file to", path.Join(workingDir, loadsOutputFilename))
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dataInfo.Lighting.Lighting); err != nil {
		fmt.Println("Error encoding lights for loads file:", err.Error())
		os.Exit(1)
	}

	if err := os.WriteFile(loadsOutputFilename, b.Bytes(), 0644); err != nil {
		fmt.Println("Error writing loads file:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Done.")

	// Write triggers
	fmt.Println("Writing trigger plist to", path.Join(workingDir, triggersOutputFilename))
	str := generateTriggers(dataInfo.Lighting.Lighting)

	if err := os.WriteFile(triggersOutputFilename, []byte(str), 0644); err != nil {
		fmt.Println("Error writing trigger plist file:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Done.")

	// Write automator services
	dir, err := getSavantDir(sysInfo, docInfo)
	if err != nil {
		fmt.Println("Error getting Savant Directory:", err.Error())
		os.Exit(1)
	}
	caches, err := getSavantCacheDir(dir)
	if err != nil {
		fmt.Println("Error getting Savant Cache Directory:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Writing workflows to", caches)

	for _, light := range dataInfo.Lighting.Lighting {
		createWorkflow(dir, caches, light)
		fmt.Println("Created Workflow for", light.Zone, light.Name)
	}
	fmt.Println("Done.")

}

func decodeFile(filepath string, data interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := plist.NewDecoder(file)
	return decoder.Decode(data)
}
