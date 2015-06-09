package goworld

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var (
	acp *Acp
)

// StartWorker start register structs and start RPC server
func StartWorker(config *Config, name string, demo bool) {
	var (
		port       int
		workerName string
	)

	for _, workerDef := range config.Workers {
		if workerDef.Name == name {
			workerName = workerDef.Name
			port = workerDef.Port
			break
		}
	}
	if port == 0 && len(workerName) == 0 {
		log.Fatalf("Error starting worker. No definition in config for name : %s\n", name)
	}
	acp = NewAcp(name)
	if !demo {
		if err := acp.Connect(workerName, 0, 1); err != nil {
			log.Panicf("ACP Connection error: %v\n", err)
		}
	}
	protocol := &Protocol{port, workerName}
	rpc.Register(protocol)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		log.Fatal("listen error:", e)
		os.Exit(1)
	}
	log.Printf("Worker serwer starts at %d\n", port)
	if err := http.Serve(l, nil); err != nil {
		log.Fatal("RPC SERVER ERROR ", err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Panic(err)
		}
	}()
}

//
// Protocol type to wrap RPC communication
//
type Protocol struct {
	Port       int
	WorkerName string
}

// GetResponse returns response object from worker
func (t *Protocol) GetResponse(request *Request, resp *Response) error {

	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()

	path := request.Path
	log.Printf("Handle path: %s\n", path)
	acp.PutString(path)
	body := []map[string]string{}
	noOfRecs := acp.GetUint()
	noOfFields := acp.GetUint()

	for i := 0; i < noOfRecs; i++ {
		m := make(map[string]string)
		for j := 0; j < noOfFields; j++ {
			fieldName := acp.GetString()
			fieldValue := acp.GetString()
			m[fieldName] = fieldValue
		}
		body = append(body, m)
	}
	resp.Err = nil
	resp.Body = body
	return nil
}

// GetResponse returns response object from worker
func (t *Protocol) GetDemoResponse(request *Request, resp *Response) error {

	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()
	path := request.Path
	log.Printf("Handle path: %s\n", path)

	body := []map[string]string{}

	m := make(map[string]string)
	m[t.WorkerName] = path
	body = append(body, m)
	resp.Err = nil
	resp.Body = body
	return nil
}
