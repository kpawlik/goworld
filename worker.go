package goworld

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
)

var (
	acp *Acp
)

const (
	SucessStatus = 0
)

// StartWorker start register structs and start RPC server
func StartWorker(config *Config, name string, mode WorkMode) {
	defer func() {
		if err := recover(); err != nil {
			log.Panic(err)
		}
	}()
	workerDef := config.GetWorkerDef(name)
	if workerDef == nil {
		log.Fatalf("Error starting worker. No definition in config for name : %s\n", name)
	}
	// start ACP
	acp = NewAcp(name)
	if mode != DemoMode {
		if err := acp.Connect(name, 0, 1); err != nil {
			log.Panicf("ACP Connection error: %v\n", err)
		}
	}
	// register worker for RPC server
	worker := &Worker{
		Port:       workerDef.Port,
		WorkerName: workerDef.Name}
	rpc.Register(worker)
	rpc.HandleHTTP()
	// start listening for requests from HTTP server
	if listener, err := net.Listen("tcp", portNo(worker.Port)); err != nil {
		log.Panicf("Start worker error on port %d. Error: %v\n", worker.Port, err)
	} else {
		log.Printf("Worker started at port: %d\n", workerDef.Port)
		log.Fatalf("RPC SERVER ERROR! %s\n", http.Serve(listener, nil))
	}
}

// Worker type to wrap RPC communication
type Worker struct {
	Port       int
	WorkerName string
	Protocol   *ProtocolConf
}

// ListObjectsFields returns response object from worker
// Demo protocol method. Returns list of fields from objects.
// All data are converted to strings
func (t *Worker) ListObjectsFields(request *Request, resp *Response) error {
	var (
		bodyElem BodyElement
	)
	defer func() {
		if r := recover(); r != nil {
			log.Panicf("PANIC in method ListObjectsFields: %v\n", r.(error))
		}
	}()
	// send protocol name to ACP
	acp.PutString(request.Protocol.Name)
	// send path
	acp.PutString(request.Path)

	status := acp.GetUbyte()
	if status != SucessStatus {
		err := acp.GetString()
		resp.Error = NewAcpErr(fmt.Sprintf("ACP ERR %v", err))
		return nil
	}
	noOfRecs := acp.GetUint()
	noOfFields := acp.GetUint()
	body := make(Body, 0, noOfRecs)
	for i := 0; i < noOfRecs; i++ {
		bodyElem = make(BodyElement)
		for j := 0; j < noOfFields; j++ {
			fieldName := acp.GetString()
			fieldValue := acp.GetString()
			bodyElem[fieldName] = fieldValue
		}
		body = append(body, bodyElem)
	}
	resp.Body = body
	return nil
}

func (t *Worker) Custom(request *Request, resp *Response) error {
	var (
		bodyElem BodyElement
	)
	protocol := request.Protocol
	path := request.Path
	pathParams := strings.Split(path, "/")
	if len(pathParams) != len(protocol.Params) {
		resp.Error = NewAcpErr("Number of params in request different then in protocol definition")
		return nil
	}
	// send protocol name to ACP
	acp.PutString(protocol.Name)
	// Send all param name and value to ACP
	for i, paramDef := range protocol.Params {
		if param, err := ParseStringParam(pathParams[i], paramDef.Type); err != nil {
			resp.Error = err
			return nil
		} else {
			if err := acp.Put(paramDef.Type, param); err != nil {
				resp.Error = NewAcpErr(fmt.Sprintf("Error puting parameter '%s' value: '%s' as type %s. Err: %v\n",
					paramDef.Name, param, paramDef.Type, err))
				return nil
			}
		}
	}
	// get status
	status := acp.GetUbyte()
	if status != SucessStatus {
		err := acp.GetString()
		resp.Error = NewAcpErr(fmt.Sprintf("Error from acp: %s", err))
		return nil
	}
	// Get Recods
	noOfRecs := acp.GetUint()
	body := make(Body, 0, noOfRecs)

	resultFieldsDef := protocol.Results
	for i := 0; i < noOfRecs; i++ {
		bodyElem = make(BodyElement)
		for _, fieldDef := range resultFieldsDef {
			if value, err := acp.Get(fieldDef.Type); err != nil {
				resp.Error = NewAcpErr(fmt.Sprintf("%v", err))
				return nil
			} else {
				bodyElem[fieldDef.Name] = value
			}
		}
		body = append(body, bodyElem)
	}
	resp.Body = body
	return nil
}
