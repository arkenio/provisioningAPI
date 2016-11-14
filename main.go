package main

import (
	"github.com/arkenio/s3provisioner/server"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	//"html"
	"net/http"
)

func main() {
	viper.AutomaticEnv()          // read in environment variables that match
	glog.Infof("ss %s", viper.GetString("ATLAS_USERNAME"));
	glog.Fatal(http.ListenAndServe(":8788", server.NewRouter()))
}
