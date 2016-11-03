package main

import (
	"github.com/arkenio/s3provisioner/server"
	"github.com/golang/glog"
	//"html"
	"net/http"
)

func main() {
	glog.Fatal(http.ListenAndServe(":8080", server.NewRouter()))
}
