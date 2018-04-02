package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wakeful/prototype-pagesnapp/browser"
	"github.com/wakeful/prototype-pagesnapp/objStore"

	"github.com/gorilla/mux"
)

const (
	NAME = "pagesnapp"
)

var (
	objStoreLocation = flag.String("store", "127.0.0.1:9090", "obj store location")
	accessKeyID      = flag.String("key", "MINIO_ACCESS_KEY_REPLACE_ME", "obj store access key")
	secretAccessKey  = flag.String("secret", "MINIO_SECRET_KEY_REPLACE_ME", "obj store secret key")

	mClient *objStore.ObjStore
)

func main() {

	flag.Parse()

	var err error
	mClient, err = objStore.NewObjStore(*objStoreLocation, *accessKeyID, *secretAccessKey, NAME)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", serveMainPageRequest)
	r.HandleFunc("/", serveMainPageRequest).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatalln(srv.ListenAndServe())

}
func serveMainPageRequest(w http.ResponseWriter, r *http.Request) {
	var html = []byte(`
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
<!-- Bootstrap CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
<title>pagesnapp</title>
</head>
<body>
<div class="container">
<div class="row">
<div class="col-sm">
<h1>pagesnapp!</h1>
<legend>Example legend</legend>
<p>add target url at the end e.q.</p>
<p><a href="?url=http://jaskula.pl">?url=http://jaskula.pl</a></p>
</div>
</div>
</div>
</body>
</html>
`)

	url := r.FormValue("url")
	if url == "" {
		w.Write(html)
		return
	}

	log.Printf("trying to take screenshot of %s", url)

	webBrowser, err := browser.NewBrowser("127.0.0.1", 4444, "firefox")
	if err != nil {
		log.Fatalln(err)
	}
	defer webBrowser.Close()

	img, err := webBrowser.TakeScreenshot(url)
	if err != nil {
		log.Fatalln(err)
	}

	tmpFile, err := ioutil.TempFile("", NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up file

	if _, err := tmpFile.Write(img); err != nil {
		log.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	name, size, err := mClient.SavePngFile(tmpFile.Name())
	if err != nil {
		log.Fatalln(err)
	}

	link, err := mClient.GenerateAccessLink(name)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write([]byte(fmt.Sprintf(""+
		"<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\">"+
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\"><title>pagesnapp</title><body>"+
		"<p>Obejct name %s<br />size %d<br />url: <a href=\"%s\">Get screenshot!</a></p></head></body></html>", name, size, link)))

}
