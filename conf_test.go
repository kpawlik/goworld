package goworld

import (
	"testing"
)

var (
	testConfig string = `
{
  "server": {
    "port": 4000,
    "protocols": [
      {
        "name": "list",
        "enabled": true
      },
      {
        "name": "search",
        "enabled": false
      },
      {
        "name": "hotel",
        "enabled": true,
        "params": [
          {
            "name": "name",
            "type": "string"
          }
        ],
        "results": [
          {
            "name": "name",
            "type": "string"
          },
          {
            "name": "address1",
            "type": "string"
          }
        ]
      }
    ]
  },
  "workers": [
    {
      "port": 4002,
      "host": "localhost",
      "name": "w1"
    },
    {
      "port": 4001,
      "host": "localhost",
      "name": "w2"
    }
  ]
}
`
)

func unm() (*Config, error) {
	b := []byte(testConfig)
	return unmarshal(b)
}

func TestUnmarshal(t *testing.T) {
	c, err := unm()
	if err != nil || c == nil {
		t.Errorf("Error parsing config: %v\n", err)
	}
}

func TestReadConfig(t *testing.T) {
	c, _ := unm()
	if c.Server.Port != 4000 {
		t.Errorf("Wrong number of server port \n")
	}
	if len(c.Workers) != 2 {
		t.Errorf("Wring number of workers")
	}
}

func TestGetProtocol(t *testing.T) {
	c, _ := unm()
	if c.GetProtocolDef("list") == nil {
		t.Errorf("Error getting protocol list")
	}
}

func TestReadProtocol(t *testing.T) {
	c, _ := unm()
	if len(c.Server.Protocols) != 3 {
		t.Errorf("Wrong number of protocols")
	}
	for _, prot := range c.Server.Protocols {
		switch prot.Name {
		case "list":
			if !prot.Enabled {
				t.Errorf("Wrong value for protocol '%s' field Enabled\n", prot.Name)
			}
		case "search":
			if prot.Enabled {
				t.Errorf("Wrong value for protocol '%s' field Enabled\n", prot.Name)
			}
		case "hotel":
			if len(prot.Params) != 1 {
				t.Errorf("Wrong no of params for protocol '%s' \n", prot.Name)
			}
			if len(prot.Results) != 2 {
				t.Errorf("Wrong no of result fields for protocol '%s' \n", prot.Name)
			}
		}
	}
}
