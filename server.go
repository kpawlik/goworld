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
	requestNames = map[WorkMode]string{
		DemoMode: "Worker.GetDemoResponse",
		TestMode: "Worker.GetTestResponse"}

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
	connectedWorkersNo := 0
	for _, workerDef := range workersDef {
		conn, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", workerDef.Host, workerDef.Port))
		if err != nil {
			worker := &WorkerConnection{Host: workerDef.Host,
				Port: workerDef.Port,
				Name: workerDef.Name}
			offlineWorkers <- worker
			continue
		}
		log.Printf("Worker %s (%s:%d), CONNECTED \n", workerDef.Name, workerDef.Host, workerDef.Port)
		connectedWorkersNo++
		worker := &WorkerConnection{Host: workerDef.Host,
			Port: workerDef.Port,
			Name: workerDef.Name,
			Conn: conn}
		onlineWorkers <- worker
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
func (t *ReqHandler) getProtocolConf(path string) *ProtocolConf {
	res := strings.Split(path, "/")
	if len(res) == 0 {
		return nil
	}
	return t.Config.GetProtocolDef(res[0])
}

// parsePath parse string parameter depend of which WorkMode has been selected.
// In case of NormalMode first part of path depends which Request method will be called.
func (t *ReqHandler) parsePath(path string) (processedPath, requestFunc string, ok bool) {
	path, _ = url.QueryUnescape(path)
	if t.WorkMode != NormalMode {
		// demo and test
		requestFunc, ok = requestNames[t.WorkMode]
		processedPath = path
		return
	}
	res := strings.Split(path, "/")
	if len(res) == 0 {
		ok = false
		return
	}
	processedPath = strings.Join(res[1:], "/")
	protocolName := res[0]
	switch protocolName {
	case "list":
		requestFunc, ok = pathToReques[protocolName]
		return
	}
	requestFunc, ok = pathToReques["custom"]
	return
}

// ServeHTTP is http request handler.
func (t *ReqHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		path, requestFuncName string
		ok                    bool
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("PANIC: %v\n", err)
		}
	}()
	path = r.URL.Path[1:]
	protocolConf := t.getProtocolConf(path)
	if protocolConf == nil || !protocolConf.Enabled {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized protocol\n")
		return
	}
	path, requestFuncName, ok = t.parsePath(path)
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Unsupported protocol")
		return
	}
	response := &Response{}
	request := &Request{
		Path:     path,
		Protocol: protocolConf}
	// get free worker
	for {
		worker := <-t.Online
		conn := worker.Conn
		if err := conn.Call(requestFuncName, request, &response); err != nil {
			log.Printf("ERROR: Response from worker %s, code: %d. Remote procedure call: %s\n",
				worker.Name, http.StatusInternalServerError, err)
			// add worker to the offline pool
			t.Offline <- worker
		} else {
			// return worker to the online pool
			t.Online <- worker
			break
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
}
