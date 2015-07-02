# goworld

Aplication which allows to access Smallworld data as JSON via HTTP

[Visit project page] (http://kpawlik.github.io/goworld)

***

## Table of Contents
- [About](#about)
  - [Features](#features)
- [Configuration](#config)
  - [Configuration file](#config-file)
  - [List protocol](#config-list-protocol)
  - [Custom protocol](#config-custom-protocol)
- [Downloads](#download)
- [Build from sources](#build)
- [Start HTTP Server and Worker](#starting)
  - [Start HTTP Server](#starting-http)
  - [Start worker](#starting-worker)
- [Tutorial](#tutorial)
  - [Quick start](#tutorial-quick-start)
  - [Define custom protocol](#tutorial-custom-protocol)
- [Protocol](#protocol)
  - [List protocol](#protocol-list)
  - [Custom protocol](#protocol-custom)
- [Stop](#stop)


***

<div id='about'/>
## About

This is a simple application which shows how data from Smallworld can be accessed via HTTP in JSON format. Goworld is composed of set of concurrent workers and one HTTP server. Worker communicates with Smallworld session via ACP protocol. HTTP server and workers communicates via RPC protocol. HTTP server publish data from workers as JSON.

<div id='features'/>
### Features 

- Zero instalation
- One executable, one Magik file
- Simple, light, fast and scalable
- Linux/Windows support
- One simple config file
- Scalable - multiple ACP workers -> one concurrency HTTP server
- You can run multiple workers on single Smallworld session
- You can run multiple workers on multiple Smallworld sessions
- HTTP Server can be run on Windows or Linux

<div id='config'/>
## Configuration

<div id='config-file'/>
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

<div id='config-list-protocol'/>
### List protocol

This is predefined protocol described in [quick start](#tutorial-quick-start).
To disable this protocol set `false` in enabled attribute, or remove whole JSON object from configuration file. 

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
   To define custom protocol you need to configure definition in configuration file.
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
 - `name` - uniqe protocol name. This will be first part  of the HTTP Request
 - `enabled` - bool attribute, which allows to disable/enable protocol
 - `params` - list of Parameter objects (name and type). Parameters values should be pass in HTTP request after protocol name and should be separated by '/' char. Parameters will be converted to appropriate type and send to ACP.

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


<div id='download'/>
## Downloads

[Binaries download](https://sourceforge.net/projects/goworld/files/?source=navbar)

[Source download](https://github.com/kpawlik/goworld)

<div id='build'/>
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

<div id='starting'/>
## Start HTTP Server and Worker

<div id='starting-http'/>
### Start HTTP Server
To start HTTP sever:
  1. Create configuration file with `Server.port`, at least list protocol enbled and at least one worker definition.
  1. Open command line terminal and type:
 
	```
	goworld.exe -t http -c [PATH TO CONFIG FILE]
	```
	HTTP server will start on defined port number, running workers from definition will be connected. Of you will start worker after you start HTTP server it will be also connected.
	
<div id='starting-worker'/>
### Start Worker
  1. To start worker you can use Magik procedure `start_goworld_worker` from file `goworld.magik`	

	```
	start_goworld_worker([NAME], [PATH TO goworld.exe], [PATH TO CONFIG FILE], [PATH TO LOG FILE])
	```
	
	`NAME` - uniqe worker name. Need to be the same as in configuration file
	
	`PATH TO goworld.exe` - path to goworld executable file
	
	`PATH TO CONFIG FILE` - path to JSON configuration file
	
	`PATH TO LOG FILE` - path where to store log file for this worker
	
  1. This procedure will start ACP process. In background it will call:

	```
	goworld.exe -n [NAME] -t worker -c [PATH TO CONFIG FILE] -l [PATH TO LOG FILE]
	```

<div id='tutorial'/>
## Tutorial

<div id='tutorial-quick-start'/>
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


<div id='tutorial-custom-protocol'/>
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

1. Create new Magik method in class goworld, this method will handle new protocol
	```
	_method goworld.find_hotel_protocol()
	## 
	## 
		!print_float_precision! << 12
		# This will get name from "params"
	    _local name << _self.get_chars()
		#send status
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

1. Register method with protocol name. Name must be the same as protocol name in config file:
	```
	goworld.register_protocol("find_hotel", :|find_hotel_protocol()|)
	$
	```
	
1. Start goworld worker:
	```
	start_goworld_worker("w1", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w1.log")
	$
	```

1. Start HTTP server
	```
	c:\tmp\goworld.exe -t http -c c:\tmp\goworld.json
	```

1. In browser type:
	```
	http://localhost:4000/find_hotel/[HOTEL_NAME]
	```

<div id='protocol'/>
## Protocols

Protocol describes how to Magik ACP and goworld worker communicate.


<div id='protocol-list'/>
### List protocol
`List protocol` starts with `list` prefix eg.

`http://localhost:4000/list/gis/hotel/100/id/name/address1`

Request structure:
`http://[HOST]:[PORT]/list/[DATASET]/[COLLECTION]/[LIMIT, 0 = ALL]/[LIST OF FIELDS SEPARATED BY "/"]``

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

1. Magik ACP send to `goworld` status code as `unsigned byte`
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

1. Magik ACP send to `goworld` number of record which will be send as `unsigned int`
	```
	_self.put_unsigned_int(records_to_get)
	```

1. Magik ACP send to `goworld` number of fields which will be send as `unsigned int`
	```
	_self.put_unsigned_int(no_of_fields)
	```

1. in the loop magik ACP sends field names and field values
	```
	_self.put_chars(field_name)
	_self.flush()
	_self.put_chars(field_value)
	_self.flush()
	```

<div id='protocol-custom'/>
### Custom protocol

To define custom prototcol you need to:

1. Define protocol in config file.
	1. Protocol name
	1. Enabled flag
	1. List of entry parameters
	1. List of result fields
1. Create magik method which will handle protocol on Smallworld side. This method need to:
    1. Receive all parameters defined in config file
	1. Send sucess code, or error code and error message
	1. Send number of records to send
	1. In loop, send fields which are defined in config file
1. Register magik method with protocol name

See example in [tutotrial](#tutorial-custom-protocol).


<div id='stop'/>
## Stop 
To stop just open system Task Manager and kill all `goworld` processes.

***


