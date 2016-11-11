package server

import (
	"encoding/json"
	s3provisioner "github.com/arkenio/s3provisioner/provisioner"
	//"github.com/golang/glog"
	"net/http"
)

type ProvisionInfo struct {
	BucketName string `json:"bucketName"`
	IamUser    string `json:"iamUser"`
	Region     string `json:"region"`
}

func ProvisionS3Post(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	provisionInfo := &ProvisionInfo{}
	err := decoder.Decode(provisionInfo)

	//TODO  should validate bucket name, IAM name and region

	bucket := &s3provisioner.S3Bucket{}
	if provisionInfo.BucketName == "" {
		provisionInfo.BucketName = s3provisioner.DEFAULT_S3IO_BUCKET_NAME
	}

	if provisionInfo.IamUser == "" {
		provisionInfo.IamUser = s3provisioner.DEFAULT_S3IO_IAMUSERNAME
	}

	if provisionInfo.Region == "" {
		provisionInfo.Region = s3provisioner.DEFAULT_AWS_S3_REGION
	}

	bucket, err = bucket.ProvisionS3AndIAMUser(provisionInfo.IamUser, provisionInfo.Region, provisionInfo.BucketName, true)

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
