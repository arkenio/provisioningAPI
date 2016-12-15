package main

import (
	"github.com/arkenio/provisioningAPI/server"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	//"html"
	"net/http"
)

func main() {
	viper.AutomaticEnv()          // read in environment variables that match
	glog.Fatal(http.ListenAndServe(":8788", server.NewRouter()))
}
