package savant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testRegistryPath = "testdata/lights.json"

func TestFromJSON(t *testing.T) {
	lights, err := fromJSON(testRegistryPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, lights)
	light := lights[0]
	assert.Equal(t, "001_1", light.ID())
	assert.Equal(t, "Upper Bath Lights", light.Name())

	assert.Equal(t, "Farah SEA.RacePointMedia_host.CurrentDimmerLevel_1_001", light.stateName())
	assert.Equal(t, "OFF", light.State())
}
