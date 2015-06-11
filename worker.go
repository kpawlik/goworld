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
func StartWorker(config *Config, name string, mode WorkMode) {
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
	if mode != DemoMode {
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

// GetDemoResponse returns response object from worker
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

// GetTestResponse returns response object from worker
func (t *Protocol) GetTestResponse(request *Request, resp *Response) error {
	var (
		testInt, resInt       int
		ok                    bool
		testFloat, resFloat   float64
		testString, resString string
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()

	body := []map[string]string{}
	bodyElem := make(map[string]string)

	path := request.Path
	log.Printf("GetTestResponse - handle path: %s\n", path)
	// test ushort
	testInt = 12
	acp.PutUshort(uint16(testInt))
	resInt = acp.GetUshort()
	ok = testInt+1 == resInt
	bodyElem["UshortTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	// test short
	testInt = -111
	acp.PutShort(int16(testInt))
	resInt = acp.GetShort()
	ok = testInt+1 == resInt
	bodyElem["ShortTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	// test uint
	testInt = 11112
	acp.PutUint(uint32(testInt))
	resInt = acp.GetUint()
	ok = testInt+1 == resInt
	bodyElem["UintTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	// test int
	testInt = -11112
	acp.PutInt(int32(testInt))
	resInt = acp.GetInt()
	ok = testInt+1 == resInt
	bodyElem["IntTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)

	//test float
	testFloat = -111112122.44
	acp.PutFloat(float64(testFloat))
	resFloat = acp.GetFloat()
	ok = testFloat+1 == resFloat
	bodyElem["FloatTest"] = fmt.Sprintf("Sent: %f, Get: %f, OK: %v", testFloat, resFloat, ok)

	//test string
	testString = "111112122.44 oóDćłą   "
	acp.PutString(testString)
	resString = acp.GetString()
	ok = testString == resString
	bodyElem["StringTest"] = fmt.Sprintf("Sent: %s, Get: %s, OK: %v", testString, resString, ok)
	body = append(body, bodyElem)
	resp.Body = body
	return nil
}
