# goworld

Test app which allows to access Smallworld data as JSON via HTTP

***

## Table of Contents
- [Quick start](#quick-start)
 - [Prerequisites](#prerequisites)
 - [Instalation](#instalation)
 - [Compilation](#compilation)
- [Example of config and run](#config)
- [Features] (#features)

***

<div id='quick-start'/>
## Quick Start

<div id='prerequisites'/>
### Prerequisites
Install Go SDK (or extract zip archive) and setup GOPATH

[Go download](https://golang.org/dl/)

[GOPATH](https://github.com/golang/go/wiki/GOPATH)

<div id='instalation'/>
### Instalation:

`go get github.com/kpawlik/goworld`
<div id='compilation'/>
### Compilation

go to `GOPATH/github.com/kpawlik/goworld/goworldc` run:

`go build main.go -o c:\tmp\goworld.exe`

***

<div id='config'/>
## Example of config and run

- Create JSON config file eg. c:\tmp\goworld.json

<pre><code>
{
	"server": {
		"port": 4000
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

- load Smallworld module from 'magik' folder to Smallworld session
- in the Smallworld console type: 
<pre><code>
start_goworld_worker("w1", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w1.log")
start_goworld_worker("w2", "c:\tmp\goworld.exe", "c:\tmp\goworld.json", "c:\tmp\w2.log")
</code></pre>

> this, will start two concurrent workers w1 and w2, which will communicate with HTTP server via rcp protocol on ports 4001 n 4002 on localhost.

- open command line and type:
<pre><code>
c:\tmp\goworld.exe -t http -c c:\tmp\goworld.json
</code></pre>

> this, will start the HTTP server on port 4000

- Run browser and type in address bar
<pre><code>
http://localhost:4000/[DATASET NAME]/[COLLECTION NAME]/[NO OF RECORDS, 0 = ALL]/[LIST OF FIELDS SEPARATED BY "/"]
eg.
http://localhost:4000/gis/hotel/100/id/name/address1/address2
</code></pre>

> application will display list 100 (or less if size of collection is less then 100) of JSON object eg.

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

***

<div id='features'/>
## Features 

- Zero instalation
- One executable 
- Linux/Windows support
- One simple config file
- Scalable - many ACP workers -> one concurrency HTTP server
- Many workers on single Smallworld session
- Handles failed requests
- ...

*Under construction*
