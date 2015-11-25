package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

var (
	session *mgo.Session
	dbName  string
)

func main() {
	var mongoServer, httpBinding string
	flag.StringVar(&mongoServer, "mongoServer", "", "The connection to the mongo server")
	flag.StringVar(&dbName, "dbName", "", "The database containing the crits collections")
	flag.StringVar(&httpBinding, "httpBinding", "", "The desired HTTP binding")
	flag.Parse()

	if mongoServer == "" || dbName == "" || httpBinding == "" {
		panic("Please set mongoServer, dbName, and httpBinding!")
	}

	var err error
	session, err = mgo.Dial(mongoServer)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	router := httprouter.New()
	router.GET("/:id", handler)
	log.Fatal(http.ListenAndServe(httpBinding, router))

}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	idLen := len(id)

	_, err := hex.DecodeString(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var constraint bson.M

	if idLen == 24 {
		// critsId
		constraint = bson.M{"_id": bson.ObjectIdHex(id)}
	} else if idLen == 32 {
		// md5
		constraint = bson.M{"md5": id}
	} else if idLen == 40 {
		// sha1
		constraint = bson.M{"sha1": id}
	} else if idLen == 64 {
		// sha256
		constraint = bson.M{"sa256": id}
	} else {
		http.NotFound(w, r)
		return
	}

	var s Sample
	session.DB(dbName).C("sample").Find(constraint).One(&s)

	if s.Id == "" {
		http.NotFound(w, r)
		return
	}

	var sc SampleChunk
	session.DB(dbName).C("sample.chunks").Find(bson.M{"files_id": bson.ObjectIdHex(s.FilesId.Hex())}).One(&sc)

	if sc.Id == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+s.FileName)
	fmt.Fprint(w, sc.Data.Data)
}
