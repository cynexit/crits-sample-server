# crits-sample-server

## Description

A very simple file server for maleware samples stored in CRITs, works without authentication or CRITs giant API overhead.

## Usage

Download the dependencies:
```bash
go get gopkg.in/mgo.v2
go get github.com/julienschmidt/httprouter
```
Start the server by using `go run`, `go build` or `go install`, whatever you prefer.
There is also an example service script for systemd if need it.

You need to provide the following flags to the service:
```
./crits-sample-server -mongoServer=1.2.3.4 -dbName=crits -httpBinding=:7889
```

You may now download samples by navigating to
```
http://localhost:7889/CRITs-ObjectId
```
The service supports download by CRITs-ObjectId, MD5, SHA1 and SHA256.