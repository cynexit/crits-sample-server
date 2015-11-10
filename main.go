package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Sample struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	FilesId  bson.ObjectId `bson:"filedata"`
	FileName string        `bson:"filename"`
}

type SampleChunk struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	N       int           `bson:"n"`
	Data    bson.Binary   `bson:"data"`
	FilesId bson.ObjectId `bson:"files_id"`
}

const (
	// your mongo server
	ServerName = "1.2.3.4"

	// the name of your crits db
	DatabaseName = "crits"

	// the http binding as "IP:PORT"
	HttpBinding = ":7889"
)

var (
	session *mgo.Session
)

func main() {
	var err error
	session, err = mgo.Dial(ServerName)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	http.HandleFunc("/", handler)
	http.ListenAndServe(HttpBinding, nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) != 24 {
		http.NotFound(w, r)
		return
	}

	var s Sample
	session.DB(DatabaseName).C("sample").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&s)

	if s.Id == "" {
		http.NotFound(w, r)
		return
	}

	var sc SampleChunk
	session.DB(DatabaseName).C("sample.chunks").Find(bson.M{"files_id": bson.ObjectIdHex(s.FilesId.Hex())}).One(&sc)

	if sc.Id == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+s.FileName)
	fmt.Fprint(w, sc.Data.Data)
}
