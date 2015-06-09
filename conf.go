package goworld

import (
	"encoding/json"
	"io/ioutil"
)

// Port - port number to listing http
type Config struct {
	Server  ServerConf
	Workers []*WorkerConf
}

type ServerConf struct {
	Port int
}

// Port no on which will be listen worker server
type WorkerConf struct {
	Host string
	Name string
	Port int
}

// read and decode JSON from file
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
