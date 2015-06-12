package goworld

import (
	"encoding/json"
	"io/ioutil"
)

// WorkMode - enumerator with serwer mode types
type WorkMode int

const (
	// Version number
	Version = "0.8"
	// UnknownMode is unrecognized mode
	UnknownMode WorkMode = iota
	// NormalMode is production mode
	NormalMode
	// DemoMode demonstration mode for test only
	DemoMode
	// TestMode  to test communication between Acp and worker
	TestMode
)

var modes = map[WorkMode]string{NormalMode: "normal", DemoMode: "demo", TestMode: "test"}

func (w WorkMode) String() string {
	return modes[w]
}

//WorkModeFromString converts string value to enumerator if not found then UnknownMode returns
func WorkModeFromString(mode string) WorkMode {
	for val, m := range modes {
		if m == mode {
			return val
		}
	}
	return UnknownMode
}

// Config application configuration structure
type Config struct {
	Server  ServerConf
	Workers []*WorkerConf
}

// ServerConf server configuration
type ServerConf struct {
	Port int
}

// WorkerConf wrkers configuration
type WorkerConf struct {
	Host string
	Name string
	Port int
}

// ReadConf reads and decodes JSON from file
func ReadConf(filePath string) (conf *Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filePath); err == nil {
		conf, err = unmarshal(data)
	}
	return
}

func unmarshal(data []byte) (conf *Config, err error) {
	err = json.Unmarshal(data, &conf)
	return
}

// Response struct
// Err - error if something go wrong
// Body - result map (field, value) to json
type Response struct {
	Err  error
	Body []map[string]string
}

//
// Request struct
// Urlreq - struct contain data to get
type Request struct {
	Path string
	Args []string
}
