package savant

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func fromJSON(filepath string) ([]*Light, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	var out []*Light
	if err := json.Unmarshal([]byte(bytes), &out); err != nil {
		return nil, err
	}

	return out, nil
}

type Light struct {
	ID             string `json:"id"`
	Zone           string `json:"zone"`
	Name           string `json:"name"`
	IsDimmer       bool   `json:"is_dimmer"`
	ReadStateName  string `json:"read_state_name"`
	WriteStateName string `json:"write_state_name"`
	Level          int    `json:"-"`
	shortName      string `json:"-"`
}

var spacersRegex = regexp.MustCompile(`[ \-_/]`)
var nonAlphaNumericRegex = regexp.MustCompile(`[^A-Za-z_\- /]`)

func apply(str string, fns ...func(string) string) string {
	for _, fn := range fns {
		str = fn(str)
	}
	return str
}

// ID is the machine-facing name for the light (eg: 001_01)
func (l Light) ShortName() string {
	if l.shortName == "" {
		l.shortName = apply(
			l.Name,
			strings.TrimSpace,
			strings.ToLower,
			func(str string) string { return nonAlphaNumericRegex.ReplaceAllString(str, "") },
			func(str string) string { return spacersRegex.ReplaceAllString(str, "_") },
		)
	}

	return l.shortName
}

// State returns the light on/off state
func (l Light) State() string {
	if l.Level > 0 {
		return "ON"
	}
	return "OFF"
}
