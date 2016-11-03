package server

import (
	"encoding/json"
	provisioner "github.com/arkenio/s3provisioner/provisioner"
	"net/http"
	"github.com/golang/glog"
)

const (
	DEFAULT_S3IO_BUCKET_NAME             = "nxios3bucket"
	DEFAULT_S3IO_IAMUSERNAME             = "nxios3bucket_iam"
	DEFAULT_S3IO_ACCESSBUCKET_POLICYNAME = "nxios3policy"
	DEFAULT_AWS_S3_REGION                = "us-east-1"
)

type ProvisionInfo struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func ProvisionS3Post(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	provisionInfo := &ProvisionInfo{}
	err := decoder.Decode(provisionInfo)

	//TODO  should validate bucket name, IAM name and region

	bucket := &provisioner.S3Bucket{}
	bucket, err = bucket.ProvisionS3AndIAMUser(provisionInfo.Name, "", provisionInfo.Name, true)
	glog.Infof("See provisioned user %v", bucket.IamUser)
	
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bucket); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func ProvisionersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
