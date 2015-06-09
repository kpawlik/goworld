package goworld

import (
	//	"fmt"
	"testing"
)

var (
	testConfig string = `
{
	"server": {
		"port": 4000
	},
	"workers": [{
		"port": 4002,
		"host": "localhost",
		"name": "w1"
	}, {
		"port": 4001,
		"host": "localhost",
		"name": "w2"
	}]
}

`
)

func TestReadConfig(t *testing.T) {
	b := []byte(testConfig)
	c := &Config{}
	c, err := unmarshal(b)
	if err != nil {
		t.Errorf("Error parsing config: %v\n", err)
	}
	if c.Server.Port != 4000 {
		t.Errorf("Wrong number of server port: %v\n", err)
	}
}
