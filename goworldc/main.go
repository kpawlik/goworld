// gosworld project main.go
package main

import (
	"flag"
	"github.com/kpawlik/goworld"
	"log"
	"os"
	"runtime"
)

//
// processName - name used to establish connection with SW ACP
// serverType - decide which server will be started
// configFilePath - path to config file
var (
	processName    string
	serverType     string
	configFilePath string
	logFile        string
	mode           string
)

// Init and parse command line params
func init() {
	var (
		file *os.File
		err  error
	)
	flag.StringVar(&processName, "n", "", "process name")
	flag.StringVar(&configFilePath, "c", "", "path to config file")
	flag.StringVar(&serverType, "t", "", "Server type [worker, http]")
	flag.StringVar(&logFile, "l", "", "logfile")
	flag.StringVar(&mode, "m", "normal", "Work mode for test and debug [normal(default), test]")
	flag.Parse()
	if serverType == "" || configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	if serverType != "http" && serverType != "worker" {
		log.Println("Wrong value in 't' switch. Allowed 'http' or 'worker' ")
		return
	}
	if serverType != "http" && processName == "" {
		log.Println("Set process name for worker")
	}
	if serverType == "http" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	if logFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		if file, err = os.Create(logFile); err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Version:     %s\n", goworld.Version)
	if len(processName) > 0 {
		log.Printf("Proces name: %s\n", processName)
	}
	log.Printf("Config file: %s\n", configFilePath)
	log.Printf("Server type: %s\n", serverType)
	if mode != "normal" {
		log.Printf("Mode type:  %s\n", mode)
	}
}

func main() {
	var (
		config   *goworld.Config
		err      error
		workMode goworld.WorkMode
	)
	if workMode = goworld.WorkModeFromString(mode); workMode == goworld.UnknownMode {
		return
	}
	if config, err = goworld.ReadConf(configFilePath); err != nil {
		log.Panicf("Error reading config file: %v\n", err)
	}

	switch serverType {
	case "http":
		startHTTPServer(config, workMode)
	case "worker":
		startWorkerServer(config, workMode)
	}
}

func startWorkerServer(config *goworld.Config, workMode goworld.WorkMode) {
	goworld.StartWorker(config, processName, workMode)
}

func startHTTPServer(config *goworld.Config, workMode goworld.WorkMode) {
	goworld.StartServer(config, workMode)
}
