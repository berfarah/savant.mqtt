package savant

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SCLIClient interface {
	Run(option string, args ...string) ([]string, error)
}

type productionSCLIClient struct {
	path string
}

var scliClient SCLIClient

const sclibridgePath = "/usr/local/bin/sclibridge"

func init() {
	isTest := strings.HasSuffix(os.Args[0], ".test")

	if !isTest {
		client, err := newSCLIClient()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		scliClient = client
	}
}

func newSCLIClient() (SCLIClient, error) {
	if _, err := os.Stat(sclibridgePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("Requires access to sclibridge")
		} else {
			return nil, err
		}
	}

	return &productionSCLIClient{path: sclibridgePath}, nil
}

func (c *productionSCLIClient) Run(option string, args ...string) ([]string, error) {
	cmdArgs := append([]string{option}, args...)
	cmd := exec.Command(sclibridgePath, cmdArgs...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(b)), "\n"), nil
}
