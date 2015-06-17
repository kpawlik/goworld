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
- One executable
- Siple, light and fast
- Linux/Windows support
- One simple config file
- Scalable - many ACP workers -> one concurrency HTTP server
- Many workers on single Smallworld session

<div id='config'/>
## Configuration

<div id='config-file'/>
### Configuration file

Config file is a simple JSON file.
<pre><code>
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
</code></pre>

`server` - HTTP server configuration

`server.port` - HTTP server port number 

`server.protocols` - list of protocols definitions. See [Protocols](#protocol)

`workers` - list of workers configuration obejcts

`worker.port` - port number to communication via RPC between HTTP server and worker

`worker.host` - host name, where worker is started

`worker.name` - unique name of worker. This will be also used as a APC process name

<div id='config-list-protocol'/>
### List protocol

This is preconfigured protocol described in [quick start](#tutorial-quick-start).
To disable this protocol set `false` in enabled attribute, or remove whole JSON object from configuration file. 
	<pre><code>
	 {
      "name": "list",
      "enabled": bool
    },
	</pre></code>

<div id='config-custom-protocol'/>
### Custom protocol
   To define custom protocol you need to configure definition in configuration file.
Configuration:
	<pre><code>
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
	</code></pre>

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


<div id='tutorial'/>
## Tutorial

<div id='tutorial-quick-start'/>
### Quick start

This example is for Windows, but this works the same way for Linux. This quick start example shows 
how to run `goworld` with example `list` protocol, which just lists fields from Smallworld objects.

- Download appropriate executable file to `C:\tmp\`

- Create JSON config file c:\tmp\goworld.json

	<pre><code>
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
	</code></pre>

- load Magik source file `goworld.magik` from `magik` folder into the Smallworld session

- in the Smallworld console type: 
	
	<pre><code>
start_goworld_worker("w1", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w1.log")
$
start_goworld_worker("w2", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w2.log")
$
	</code></pre>

   this, will start two concurrent workers `w1` and `w2`, which will communicate with HTTP server via RPC protocol on ports 4001 n 4002.
Procedure `start_goworld_worker` takes 4 parameters: 
   * `worker_name`
   * `path to goworld.exe`
   * `path to config file`
   * `path where worker log file will be written`

- open Windows command line and type:
<pre><code>
c:\tmp\goworld.exe -t http -c c:\tmp\goworld.json
</code></pre>

   this, will start the HTTP server on port 4000

- start internet browser and type in address bar
	<pre><code>
http://localhost:4000/list/[DATASET NAME]/[COLLECTION NAME]/[LIMIT, 0 = ALL]/[LIST OF FIELDS SEPARATED BY "/"]
eg.
http://localhost:4000/list/gis/hotel/100/id/name/address1/address2
	</code></pre>

   application will display list 100 (or less if size of collection is less then 100) of JSON object eg.

	<pre><code>
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
	</code></pre>


<div id='tutorial-custom-protocol'/>
### Define custom protocol

1. Get config file from [quick start](#tutorial-quick-start). Add new protocol definition with name `find_hotel`
	
	<pre><code>
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
	</code></pre>

1. Create new Magik method in class goworld, this method will handle new protocol

	<pre><code>
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
	</code></pre>

1. Register method with protocol name. Name must be the same as in config file:
	
	<pre><code>
	goworld.register_protocol("find_hotel", :|find_hotel_protocol()|)
	$

	
1. Start goworld worker:
	<pre><code>
	start_goworld_worker("w1", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w1.log")
	$
	</code></pre>

1. Start HTTP server
	<pre><code>
	c:\tmp\goworld.exe -t http -c c:\tmp\goworld.json
	</code></pre>

1. In browser type:
	<pre><code>
	http://localhost:4000/find_hotel/[HOTEL_NAME]
	</code></pre>

<div id='protocol'/>
## Protocols

Prototcol means the way how magik Acp and goworld worker are communicate.


<div id='protocol-list'/>
### List protocol
List protocol starts with `list` prefix eg.

`http://localhost:4000/list/gis/hotel/100/id/name/address1`

Request structure:
`http://[HOST]:[PORT]/list/[DATASET]/[COLLECTION]/[LIMIT, 0 = ALL]/[LIST OF FIELDS SEPARATED BY "/"]``

Protocol returns a list of JSON objects, each JSON object contain LIST OF FIELDS and VALUES from COLLECTION. Number of objects is limited by LIMIT or size of COLLECTION.
If DATASET or COLLECTION doesn not exists, error message is returned. If field does not exists error is returned as a field value.

Communication:

1. `goworld` send to magik a char vector, this is a path from http address. eg for 
`http://localhost:4000/list/gis/hotel/100/id/name/address1/address2` it will be `gis/hotel/100/id/name/address1/address2`.

	<pre><code>
	_local path << _self.get_chars()
	</code></pre>

1. magik ACP send to goworld information about status as `unsigned byte`
   * `0` means `OK`, no error

	<pre><code>
	_self.put_unsigned_byte(0)
	</code></pre>
   * `> 0` means error (no such dataset, no access). In this case as a next message magik need to send string with error message
	<pre><code>
	_self.put_unsigned_byte(1)
	_self.chars(write_string("No such dataset with name", dataset_name))
	_continue
	</code></pre>


1. magik ACP send to goworld number of record which will be send as `unsigned int`
	<pre><code>
	_self.put_unsigned_int(records_to_get)
	</code></pre>

1. magik ACP send to goworld number of fields which will be send as `unsigned int`
	<pre><code>
	_self.put_unsigned_int(no_of_fields)
	</code></pre>

1. in the loop magik ACP sends field names and field values
	<pre><code>
	_self.put_chars(field_name)
	_self.flush()
	_self.put_chars(field_value)
	_self.flush()
	</code></pre>

<div id='protocol-custom'/>
### Custom protocol

To define custom prototcol you need to:

1. Define protocol in config file.
	1. Protocol name
	1. enabled flag
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


