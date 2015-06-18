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
	// SucessStatus value which should be returned from ACP if no error ocure
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

	if err := acp.Connect(name, 0, 1); err != nil {
		log.Panicf("ACP Connection error: %v\n", err)
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
func (w *Worker) ListObjectsFields(request *Request, resp *Response) error {
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
	// get status
	if err := w.checkAcpStatus(); err != nil {
		resp.Error = err
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

// Custom handles communication defined by custom protocol in config file
func (w *Worker) Custom(request *Request, resp *Response) (err error) {
	var (
		bodyElem BodyElement
		acpErr   *AcpErr
	)
	protocol := request.Protocol
	pathParams := strings.Split(request.Path, "/")
	if len(pathParams) != len(protocol.Params) {
		resp.Error = NewAcpErr("Wrong number of parameters.")
		return nil
	}
	// send protocol name to ACP
	acp.PutString(protocol.Name)
	// Send all param name and value to ACP
	if acpErr = w.sendParameters(protocol, pathParams); acpErr != nil {
		resp.Error = acpErr
		return
	}
	// get status
	if acpErr = w.checkAcpStatus(); acpErr != nil {
		resp.Error = acpErr
		return
	}
	// Get Recods
	noOfRecs := acp.GetUint()
	body := make(Body, 0, noOfRecs)

	resultFieldsDef := protocol.Results
	for i := 0; i < noOfRecs; i++ {
		if bodyElem, acpErr = w.getFields(resultFieldsDef); acpErr != nil {
			resp.Error = acpErr
			return
		}
		body = append(body, bodyElem)
	}
	resp.Body = body
	return
}

// chackAcpStatus checks if ACP returns valid sucess status. If not then read error message and create error object
func (w *Worker) checkAcpStatus() *AcpErr {
	status := acp.GetUbyte()
	if status != SucessStatus {
		return NewAcpErr(fmt.Sprintf("Status error from ACP: Code %d. Message: %s", status, acp.GetString()))
	}
	return nil
}

//sendParameters sends list of parameters from request to ACP
func (w *Worker) sendParameters(protocol *ProtocolConf, pathParams []string) (err *AcpErr) {
	for i, paramDef := range protocol.Params {
		if err = w.sendParameter(paramDef, pathParams[i]); err != nil {
			return
		}
	}
	return
}

//sendParameter convert and send string parameter to ACP
func (w *Worker) sendParameter(paramDef *ParameterConf, value string) (acpErr *AcpErr) {
	var (
		param interface{}
		err   error
	)
	if param, err = ParseStringParam(value, paramDef.Type); err != nil {
		acpErr = NewAcpErr(fmt.Sprint(err))
		return
	}
	if err = acp.Put(paramDef.Type, param); err != nil {
		acpErr = NewAcpErr(fmt.Sprint(err))
		return
	}
	return
}

// getFields reads list of fields from ACP to map
func (w *Worker) getFields(resultFieldsDef []*ParameterConf) (bodyElem BodyElement, acpErr *AcpErr) {
	var (
		value interface{}
	)
	bodyElem = make(BodyElement, len(resultFieldsDef))
	for _, fieldDef := range resultFieldsDef {
		if value, acpErr = acp.Get(fieldDef.Type); acpErr != nil {
			return
		}
		bodyElem[fieldDef.Name] = value
	}
	return
}
