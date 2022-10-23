package savant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testRegistryPath = "testdata/lights.json"

func TestFromJSON(t *testing.T) {
	lights, err := fromJSON(testRegistryPath)
	assert.NoError(t, err)
	assert.Len(t, lights, 3)
	lightOne := lights[0]
	lightTwo := lights[1]
	lightThree := lights[2]

	assert.Equal(t, "001_1", lightOne.ID)
	assert.Equal(t, "Master Lights", lightOne.Name)
	assert.Equal(t, "Farah SEA.RacePointMedia_host.CurrentDimmerLevel_1_001", lightOne.ReadStateName)
	assert.Equal(t, "userDefined.SetDimmerLevel_1_001", lightOne.WriteStateName)
	assert.Equal(t, true, lightOne.IsDimmer)
	assert.Equal(t, "master_lights", lightOne.ShortName())
	assert.Equal(t, "OFF", lightOne.State())

	assert.Equal(t, "002_1", lightTwo.ID)
	assert.Equal(t, "Master Nancy's Nightstand", lightTwo.Name)
	assert.Equal(t, "Farah SEA.RacePointMedia_host.CurrentDimmerLevel_1_002", lightTwo.ReadStateName)
	assert.Equal(t, "userDefined.SetDimmerLevel_1_002", lightTwo.WriteStateName)
	assert.Equal(t, true, lightTwo.IsDimmer)
	assert.Equal(t, "master_nancys_nightstand", lightTwo.ShortName())
	assert.Equal(t, "OFF", lightTwo.State())

	assert.Equal(t, "003_1", lightThree.ID)
	assert.Equal(t, "Lower Bath Closet", lightThree.Name)
	assert.Equal(t, "Farah SEA.RacePointMedia_host.CurrentDimmerLevel_1_003", lightThree.ReadStateName)
	assert.Equal(t, "userDefined.SetDimmerLevel_1_003", lightThree.WriteStateName)
	assert.Equal(t, false, lightThree.IsDimmer)
	assert.Equal(t, "lower_bath_closet", lightThree.ShortName())
	assert.Equal(t, "OFF", lightThree.State())
}
