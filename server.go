package goworld

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"strings"
	"time"
)

var (
	requestNames = map[WorkMode]string{
		DemoMode: "Protocol.GetDemoResponse",
		TestMode: "Protocol.GetTestResponse"}

	pathToReques = map[string]string{"list": "Protocol.ListObjectsFields"}
)

// StartServer initialize workers and starts HTTP server
func StartServer(config *Config, mode WorkMode) {
	online, offline := initWorkersConnections(config.Workers)
	port := fmt.Sprintf(":%d", config.Server.Port)
	// try reconect workers in goroutine
	go handleOfflineWorkers(online, offline)
	startHTTPServer(port, online, offline, mode)
	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()
}

// startHTTPServer starts main http server
func startHTTPServer(port string, online chan *Worker, offline chan *Worker, mode WorkMode) {
	log.Printf("Http server start on port: %s\n", port)
	server := &http.Server{
		Addr: port,
		Handler: &Handler{Online: online,
			Offline:  offline,
			WorkMode: mode},
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("SERVWER ERROR ! %s\n", server.ListenAndServe())
}

// initWorkersConnections initialize RPC clients connections. Cache them to online channel.
// Workers which are not connected are send to offline channel
func initWorkersConnections(workersDef []*WorkerConf) (onlineWorkers chan *Worker, offlineWorkers chan *Worker) {
	workersNo := len(workersDef)
	onlineWorkers = make(chan *Worker, workersNo)
	offlineWorkers = make(chan *Worker, workersNo)
	connectedWorkersNo := 0
	for _, workerDef := range workersDef {
		conn, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", workerDef.Host, workerDef.Port))
		if err != nil {
			worker := &Worker{Host: workerDef.Host,
				Port: workerDef.Port,
				Name: workerDef.Name}
			offlineWorkers <- worker
			continue
		}
		log.Printf("Worker %s (%s:%d), CONNECTED \n", workerDef.Name, workerDef.Host, workerDef.Port)
		connectedWorkersNo++
		worker := &Worker{Host: workerDef.Host,
			Port: workerDef.Port,
			Name: workerDef.Name,
			Conn: conn}
		onlineWorkers <- worker
	}
	log.Printf("Connected workers: %d \n", connectedWorkersNo)
	return
}

//handleOfflineWorkers trying to reconnect offline workers every one second. When worker will reconnect
// send him to online chanel
func handleOfflineWorkers(online chan *Worker, offline chan *Worker) {
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

// Worker type to store worker conenction and data
type Worker struct {
	Name string
	Host string
	Port int
	Conn *rpc.Client
}

// Handler struct to implement ServeHTTP method
type Handler struct {
	Online   chan *Worker
	Offline  chan *Worker
	WorkMode WorkMode
}

// parsePath parse path depend of which WorkMode has been selected.
// In case of NormalMode first part of path  depends which Request method will be called.
func (t *Handler) parsePath(path string) (processedPath, requestName string, ok bool) {
	if t.WorkMode != NormalMode {
		requestName, ok = requestNames[t.WorkMode]
		processedPath = path
		return
	}
	res := strings.Split(path, "/")
	if len(res) < 1 {
		ok = false
		return
	}
	pathPrefix := res[0]
	requestName, ok = pathToReques[pathPrefix]
	processedPath = strings.Join(res[1:], "/")
	return
}

// ServeHTTP is http request handler.
func (t *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		path, requestName string
		ok                bool
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("PANIC: %v\n", err)
		}
	}()
	path = r.URL.Path[1:]
	path, requestName, ok = t.parsePath(path)
	if !ok {
		log.Printf("Unsuported protocol %s\n", path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Unsupported protocol")
		return
	}
	log.Printf("Handle request with path: %s\n", path)
	response := &Response{}
	request := &Request{Path: path}
	// get free worker
	for {
		worker := <-t.Online
		log.Printf("Send request %s to worker %s\n", path, worker.Name)
		conn := worker.Conn
		if err := conn.Call(requestName, request, &response); err != nil {
			// in case of error
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("ERROR: Response from worker %s, code: %d. Remote procedure call: %s\n",
				worker.Name, http.StatusInternalServerError, err)
			t.Offline <- worker
		} else {
			// return worker to the pool
			t.Online <- worker
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(response.Body); err != nil {
		log.Printf("Error Encode response body: %s\n", err)
	}

}
