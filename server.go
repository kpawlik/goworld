package goworld

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"net/url"
	"strings"
	"time"
)

var (
	pathToReques = map[string]string{
		"list":   "Worker.ListObjectsFields",
		"custom": "Worker.Custom",
	}
)

// channels with workers connections
type workerChan chan *WorkerConnection

// StartServer initialize workers and starts HTTP server
func StartServer(config *Config, mode WorkMode) {
	online, offline := initWorkersConnections(config.Workers)

	// try reconect workers in goroutine
	go handleOfflineWorkers(online, offline)
	startHTTPServer(config, online, offline, mode)
	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()
}

// startHTTPServer starts main http server
func startHTTPServer(config *Config, online workerChan, offline workerChan, mode WorkMode) {
	port := config.Server.Port
	log.Printf("HTTP server started on port: %v\n", port)
	server := &http.Server{
		Addr: portNo(port),
		Handler: &ReqHandler{Online: online,
			Offline:  offline,
			Config:   config,
			WorkMode: mode},
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("SERVWER ERROR!\n %s\n", server.ListenAndServe())
}

// initWorkersConnections initialize RPC clients connections. Cache them to online channel.
// Workers which are not connected are send to offline channel
func initWorkersConnections(workersDef []*WorkerConf) (onlineWorkers workerChan, offlineWorkers workerChan) {
	workersNo := len(workersDef)
	onlineWorkers = make(workerChan, workersNo)
	offlineWorkers = make(workerChan, workersNo)

	for _, workerDef := range workersDef {
		conn, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", workerDef.Host, workerDef.Port))
		if err != nil {
			// add to online pool
			offlineWorkers <- &WorkerConnection{Host: workerDef.Host,
				Port: workerDef.Port,
				Name: workerDef.Name}
			continue
		}
		// add to offline pool
		onlineWorkers <- &WorkerConnection{Host: workerDef.Host,
			Port: workerDef.Port,
			Name: workerDef.Name,
			Conn: conn}
		log.Printf("Worker %s (%s:%d), CONNECTED \n", workerDef.Name, workerDef.Host, workerDef.Port)
	}
	return
}

//handleOfflineWorkers trying to reconnect offline workers every one second. When worker will reconnect
// send him to online chanel
func handleOfflineWorkers(online workerChan, offline workerChan) {
	for {
		<-time.After(1 * time.Second)
		for worker := range offline {
			conn, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", worker.Host, worker.Port))
			if err != nil {
				offline <- worker
			} else {
				log.Printf("Worker %s (%s:%d), CONNECTED \n", worker.Name, worker.Host, worker.Port)
				worker.Conn = conn
				online <- worker
			}
		}
	}
}

// WorkerConnection type to store worker conenction and data
type WorkerConnection struct {
	Name string
	Host string
	Port int
	Conn *rpc.Client
}

// ReqHandler struct to implement ServeHTTP method
type ReqHandler struct {
	Online   workerChan
	Offline  workerChan
	Config   *Config
	WorkMode WorkMode
}

// getProtocolConf return protocol definition with name which should be
// first part of path
func (r *ReqHandler) getProtocolConf(protocolName string) *ProtocolConf {
	return r.Config.GetProtocolDef(protocolName)
}

// getRequestFunctionName return name of method which will handle RCP request.
// Second return value is bool, which indicates if function name was found or not.
func (r *ReqHandler) getRequestFunctionName(protocolName string) (string, bool) {
	if r.WorkMode == TestMode {
		return "Worker.GetTestResponse", true
	}
	if protocolName != "list" {
		protocolName = "custom"
	}
	reqFucn, ok := pathToReques[protocolName]
	return reqFucn, ok

}

// checkProtocolConf returns true if protocol is valid and enabled
func (r *ReqHandler) checkProtocolConf(protocolConf *ProtocolConf) bool {
	return protocolConf != nil && protocolConf.Enabled

}

//writeErrorStatus writes error on writter, error message depends on status value
func (r *ReqHandler) writeErrorStatus(w http.ResponseWriter, status int) {
	var message string
	switch status {
	case http.StatusMethodNotAllowed:
		message = "Unsupported protocol"
	case http.StatusUnauthorized:
		message = "Unauthorized protocol"
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, message)
}

// parsePath parse string parameter depend of which WorkMode has been selected.
// In case of NormalMode first part of path depends which Request method will be called.
func (r *ReqHandler) parsePath(path string) (protocolName, processedPath string, ok bool) {
	var (
		res []string
		err error
	)
	if path, err = url.QueryUnescape(path); err != nil {
		log.Printf("Error unescape path: %v\n", err)
	}
	res = strings.Split(path, "/")
	if ok = len(res) > 0; ok {
		processedPath = strings.Join(res[1:], "/")
		protocolName = res[0]
	}
	return
}

// ServeHTTP is http request handler.
func (r *ReqHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		path, requestFuncName string
		protocolName          string
		ok                    bool
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("PANIC: %v\n", err)
		}
	}()
	start := time.Now()
	path = req.URL.Path[1:]
	// check path
	if protocolName, path, ok = r.parsePath(path); !ok {
		r.writeErrorStatus(w, http.StatusMethodNotAllowed)
		return
	}
	// check request function name
	if requestFuncName, ok = r.getRequestFunctionName(protocolName); !ok {
		r.writeErrorStatus(w, http.StatusMethodNotAllowed)
		return
	}
	protocolConf := r.getProtocolConf(protocolName)
	if !r.checkProtocolConf(protocolConf) {
		r.writeErrorStatus(w, http.StatusUnauthorized)
		return
	}
	response := &Response{}
	request := &Request{
		Path:     path,
		Protocol: protocolConf}
	for {
		// get free worker from online pool
		worker := <-r.Online
		conn := worker.Conn
		if err := conn.Call(requestFuncName, request, &response); err == nil {
			// return worker to the online pool
			r.Online <- worker
			break
		} else {
			log.Printf("ERROR: Response from worker %s, code: %d. Remote procedure call: %s\n",
				worker.Name, http.StatusInternalServerError, err)
			// add worker to the offline pool and get next worker from online pool. dont break
			r.Offline <- worker
		}
	}
	if response.Error != nil {
		fmt.Fprintf(w, "%v", response.Error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(response.Body); err != nil {
		log.Printf("Error Encode response body: %s\n", err)
	}
	log.Printf("Request. Protocol: '%s', Params: '%s', Processed in %v\n", protocolName, path, time.Now().Sub(start))
}
