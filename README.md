# crits-sample-server

## Description

A very simple file server for maleware samples stored in CRITs, works without authentication or CRITs giant API overhead.

## Usage

Download the only dependencie, mgo:
```bash
go get gopkg.in/mgo.v2
```
Edit `main.go` and fill in the constants with your own values.
```go
const (

	// your mongo server
	ServerName = "1.2.3.4"

	// the name of your crits db
	DatabaseName = "crits"

	// the http binding as "IP:PORT"
	HttpBinding = ":7889"
)
```
Start the server by using `go run`, `go build` or `go install`, whatever you prefer.
There is also an example service script for systemd if need it.

You may now download samples by navigating to
```
http://localhost:7889/?id=CRITs-ObjectId
```
Of course you need to replace localhost and 7889 with your settings. Also you have to replace `CRITs-ObjectId` with an actuall CRITs sample id.