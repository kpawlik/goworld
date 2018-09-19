Table of Contents                                                    
=================                                                    

* [goworld](#goworld)                                                
  * [About](#about)                                                  
    * [Features](#features)                                          
  * [Configuration](#configuration)                                  
    * [Configuration file](#configuration-file)                      
    * [List protocol](#list-protocol)                                
  * [Download binaries and source](#download-binaries-and-source)    
  * [Build from sources](#build-from-sources)                        
    * [Prerequisites](#prerequisites)                                
    * [Instalation](#instalation)                                    
    * [Compilation](#compilation)                                    
  * [Start HTTP Server and Worker](#start-http-server-and-worker)    
    * [Start HTTP Server](#start-http-server)                        
    * [Start Worker](#start-worker)                                  
  * [Tutorial](#tutorial)                                            
    * [Quick start](#quick-start)                                    
    * [Define custom protocol](#define-custom-protocol)              
  * [Protocols](#protocols)                                          
    * [List protocol](#list-protocol-1)                              
    * [Custom protocol](#custom-protocol)                            
  * [Stop](#stop)                                                    

# goworld

Aplication which allows to access Smallworld data as JSON via HTTP

[Visit project page] (http://kpawlik.github.io/goworld)


## About

This is an application to get data from Smallworld via HTTP in JSON format. Goworld is composed of set of concurrent workers and one HTTP server. 
Worker communicates with Smallworld session via ACP protocol. 
HTTP server and workers communicates via RPC protocol. 


### Features 

- Zero installation
- One executable file, one Magik file
- Simple, light, fast and scalable
- Linux/Windows support
- One simple config file
- Scalable - multiple ACP workers -> one concurrency HTTP server
- You can run multiple workers on single Smallworld session
- You can run multiple workers on multiple Smallworld sessions
- HTTP Server can be run on Windows or Linux


## Configuration

### Configuration file

Config file is a simple JSON file.

```
{
    "server": {
        "port": Number,
		"protocols": 
			[
			...
			]
    },
    "workers": [
	    {
        	"port": Number,
	        "host": String,
	        "name": String
    	},
	...

	]
}
```

`server` - HTTP server configuration

`server.port` - HTTP server port number 

`server.protocols` - list of protocols definitions. See [Protocols](#protocol)

`workers` - list of workers definitions

`worker.port` - port number to communicate with HTTP server

`worker.host` - name of the host where worker is started

`worker.name` - unique name of the worker. This name will be also used as a Magik ACP process name


### List protocol

This is predefined protocol described in [quick start](#tutorial-quick-start).
To disable this protocol set `false` for attribute `enabled`, or just remove JSON object from configuration file. 

	```
	protocols:
	[
		{
		  		"name": "list",
		  		"enabled": bool
			},
	...
	]
	```

<div id='config-custom-protocol'/>
### Custom protocol
   To define custom protocol you need to add protocol definition in configuration file.
Configuration:

	```
	 {
	  "name": string,
	  "enabled": bool,
	  "params": [
	    {
	      "name": string,
	      "type": string
	    }
		...
	  ],
	  "results": [
	    {
	      "name": string,
	      "type": string
	    }
		...
	  ]
	},
	```
Definiton of custom protocol contains:
 - `name` - uniqe protocol name. This will be first part of the HTTP Request
 - `enabled` - bool attribute, which allows to disable/enable protocol
 - `params` - list of Parameter objects (name and type). Parameters values must be pass in request URL after protocol name and should be separated by '/' char. Parameters will be converted to appropriate type and send to ACP.

  ```
  eg. 
  http://localhost:4000/protocol_name/param1/param2/param3
  ```
 - `results` - list of fields names and types which will be received from ACP.

Supported types:
 - boolean
 - unsigned_byte
 - signed_byte
 - unsigned_short
 - signed_short
 - unsigned_int
 - signed_int
 - unsigned_long
 - signed_long
 - short_float
 - float
 - chars



## Download binaries and source

[Binaries download](https://github.com/kpawlik/goworld/releases)

[Source download](https://github.com/kpawlik/goworld)


## Build from sources

### Prerequisites
Install Go SDK (or extract zip archive) and setup GOPATH

[Go download](https://golang.org/dl/)

[GOPATH](https://github.com/golang/go/wiki/GOPATH)


### Instalation

`go get github.com/kpawlik/goworld`

### Compilation

go to `GOPATH/github.com/kpawlik/goworld/goworldc` run:

`go build main.go -o c:\tmp\goworld.exe`


## Start HTTP Server and Worker


### Start HTTP Server
To start HTTP sever:
  1. Create configuration file with `Server.port`, at least `list` protocol enbled and at least one worker definition.
  2. Open command line terminal and type:

  ```
  goworld.exe -t http -c [PATH TO CONFIG FILE]
  ```
  HTTP server will start on defined port number, running workers from definition will be connected. Workers can be start before or after you start HTTP.


### Start Worker
  1. To start worker you can use Magik procedure `start_goworld_worker` from file `goworld.magik`	

  ```
  start_goworld_worker([NAME], [PATH TO goworld.exe], [PATH TO CONFIG FILE], [PATH TO LOG FILE])
  ```

  `NAME` - uniqe worker name. Need to be the same as in configuration file

  `PATH TO goworld.exe` - path to goworld executable file

  `PATH TO CONFIG FILE` - path to JSON configuration file

  `PATH TO LOG FILE` - path where to store log file for this worker

  2. This procedure will start ACP process. In background it will call:

  ```
  goworld.exe -n [NAME] -t worker -c [PATH TO CONFIG FILE] -l [PATH TO LOG FILE]
  ```


## Tutorial


### Quick start

This example is for Windows, but this works the same way for Linux. This quick start example shows 
how to run `goworld` with example `list` protocol, which just lists fields from Smallworld objects.

- Download appropriate executable file to `C:\tmp\`

- Create JSON config file `C:\tmp\goworld.json`

  ```
  {
  	"server": {
  		"port": 4000,
  		"protocols": 
  			[
  				{
  					"name": "list",
  					"enabled": true
  				}
  			]
  	},
  	"workers": [{
  		"port": 4002,
  		"host": "localhost",
  		"name": "w1"
  	}, {
  		"port": 4001,
  		"host": "localhost",
  		"name": "w2"
  	}]
  }
  ```

- load Magik source file `goworld.magik` from `magik` folder into the Smallworld session

- in the Smallworld console type: 

  ```
  start_goworld_worker("w1", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w1.log")
  $
  start_goworld_worker("w2", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w2.log")
  $
  ```

   this, will start two concurrent workers `w1` and `w2`, which will communicate with HTTP server via RPC protocol on ports 4001 n 4002.
  Procedure `start_goworld_worker` takes 4 parameters: 
   * `worker_name` - unique name which need to be same as in  configuration file
   * `path to goworld.exe`
   * `path to config file`
   * `path where worker log file will be written`

- open Windows command line and type:
  ```
  C:\tmp\goworld.exe -t http -c C:\tmp\goworld.json
  ```

   this, will start the HTTP server on port 4000

- start internet browser and type in address bar
  ```
  http://localhost:4000/list/[DATASET NAME]/[COLLECTION NAME]/[LIMIT, 0 = ALL]/[LIST OF FIELDS SEPARATED BY "/"]
  eg.
  http://localhost:4000/list/gis/hotel/100/id/name/address1/address2
  ```

   application will display list 100 (or less if size of collection is less then 100) of JSON object eg.

  ```
  [
  	{
  		"address1": "154 Palad Road",
  		"address2": "Wyoming",
  		"id": "condition(does_not_understand: Object hotel133:(AAA Guest House) does not understand message id)",
  		"name": "AAA Guest House"
  	},
  	{
  		"address1": "219 Pinguin Road",
  		"address2": "Miami",
  		"id": "condition(does_not_understand: Object hotel133:(All 5 Seasons) does not understand message id)",
  		"name": "All 5 Seasons"
  	},
  	...
  	...
  	...
  ]
  ```



### Define custom protocol

1. Create config file same as is described in [quick start](#tutorial-quick-start). Add new protocol definition with name `find_hotel`
  ```
  ...
  "protocols": [
  	{
        	"name": "find_hotel",
          "enabled": true,
          "params": [
              {
                  "name": "name",
                  "type": "chars"
              }
          ],
          "results": [
              {
                  "name": "address_1",
                  "type": "chars"
              },
              {
                  "name": "address_2",
                  "type": "chars"
              },
              {
                  "name": "x",
                  "type": "float"
              },
  			{
                  "name": "y",
                  "type": "float"
              }
       	 	]
    	}
  ],
  ...
  ```

2. Create new Magik method in class goworld, this method will handle new protocol
  ```ruby
  _method goworld.find_hotel_protocol()
  ## 
  ## 
  	!print_float_precision! << 12
  	# This will get name from "params"
      	_local name << _self.get_chars()
  	# send status
  	_self.send_success_status()
  	_self.flush()
  	_local ds << gis_program_manager.databases[:gis]
  	_local coll << ds.collections[:hotel]
  	_local select << coll.select(predicate.eq(:name, name))
  	# send no of recs
  	_self.put_unsigned_int(select.size)
  	_self.flush()       
  	# send results fields in the same order as in config file 
  	_for rec _over select.fast_elements()
  	_loop 
  	    _self.put_chars(write_string(rec.address1))
  	    _self.flush()
  	    _self.put_chars(write_string(rec.address2))
  	    _self.flush()
  	    _self.put_float(rec.location.x)
  	    _self.flush()
  	    _self.put_float(rec.location.y)
  	    _self.flush()
  	 _endloop 
  _endmethod
  $
  ```

3. Register method with protocol name. Name must be the same as protocol name in config file:
  ```
  goworld.register_protocol("find_hotel", :|find_hotel_protocol()|)
  $
  ```

4. Start goworld worker:
  ```
  start_goworld_worker("w1", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w1.log")
  $
  ```

5. Start HTTP server
  ```
  c:\tmp\goworld.exe -t http -c c:\tmp\goworld.json
  ```

6. In browser type:
  ```
  http://localhost:4000/find_hotel/[HOTEL_NAME]
  ```


## Protocols

Protocol describes how to Magik ACP and goworld worker communicate.

### List protocol
`List protocol` starts with `list` prefix eg.

`http://localhost:4000/list/gis/hotel/100/id/name/address1`

Request structure:
`http://[HOST]:[PORT]/list/[DATASET]/[COLLECTION]/[LIMIT, 0 = ALL]/[LIST OF FIELDS SEPARATED BY "/"]` 

Protocol returns a list of JSON objects, each JSON object contain LIST OF FIELDS and VALUES from COLLECTION. Number of objects is limited by LIMIT or size of COLLECTION.
If DATASET or COLLECTION does not exists, error message will be returned. 
If field with requested name does not exists, then error will be returned as a field value.

Communication:

1. `goworld` send to Magik ACP a vector of chars, this is a part of HTTP address which occurs after protocol name. 
  `http://localhost:4000/list/gis/hotel/100/id/name/address1/address2` it will be `gis/hotel/100/id/name/address1/address2`.
  Magik code:
  ```
  _local path << _self.get_chars()
  ```

2. Magik ACP send to `goworld` status code as `unsigned byte`
   * `0` means no error
   ```
   _self.put_unsigned_byte(0)
   ```
   * `> 0` means error (no such dataset, no access etc.). In this case as a next Magik ACP needs to send string with error message
   ```
   _self.put_unsigned_byte(1)
   _self.chars(write_string("No such dataset with name", dataset_name))
   _continue
   ```

3. Magik ACP send to `goworld` number of record which will be send as `unsigned int`
  ```
  _self.put_unsigned_int(records_to_get)
  ```

4. Magik ACP send to `goworld` number of fields which will be send as `unsigned int`
  ```
  _self.put_unsigned_int(no_of_fields)
  ```

5. in the loop magik ACP sends field names and field values
  ```
  _self.put_chars(field_name)
  _self.flush()
  _self.put_chars(field_value)
  _self.flush()
  ```


### Custom protocol

To define custom protocol you need to:

1. Define protocol in config file.
  1. Protocol name
  2. Enabled flag
  3. List of entry parameters
  4. List of result fields
2. Create magik method which will handle protocol on Smallworld side. This method need to:
    1. Receive all parameters defined in config file
    2. Send sucess code, or error code and error message
    3. Send number of records to send
    4. In loop, send fields which are defined in config file
3. Register magik method with protocol name

See example in [tutotrial](#tutorial-custom-protocol).


## Stop 
To stop just open system Task Manager and kill all `goworld` processes.

***
