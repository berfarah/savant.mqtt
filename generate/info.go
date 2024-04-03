package main

const SystemInfoFilename = "rpmSystemInfo.plist"

type SystemInfo struct {
	SystemName     string `plist:"RPMSystemName"`        // eg: Farah SEA
	DeviceClassKey string `plist:"RPMUICDeviceClassKey"` // eg: RacePointMedia_host
}

const DocumentInfoFilename = "documentInfo.plist"

type DocumentInfo struct {
	VersionKey  string `plist:"RPMAppDocumentBlueprintBuildNumberKey"` // eg: 10.5.1
	BuildNumber string `plist:"RPMAppDocumentBlueprintVersionKey"`     // eg: 532
}
