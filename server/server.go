package server

import (
	"encoding/json"
	provisioner "github.com/arkenio/s3provisioner/provisioner"
	//"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
)

func ProvisionS3Post(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	provisionInfo := &provisioner.S3ProvisionInfo{}
	err := decoder.Decode(provisionInfo)
	//TODO  should validate bucket name, IAM name and region

	bucket := &provisioner.S3Bucket{}
	if provisionInfo.BucketName == "" {
		provisionInfo.BucketName = provisioner.DEFAULT_S3IO_BUCKET_NAME
	}

	if provisionInfo.IamUser == "" {
		provisionInfo.IamUser = provisioner.DEFAULT_S3IO_IAMUSERNAME
	}

	if provisionInfo.Region == "" {
		provisionInfo.Region = provisioner.DEFAULT_AWS_S3_REGION
	}
	
	w.Header().Add("Content-Type", "application/json")
	bucket, err = bucket.ProvisionS3AndIAMUser(provisionInfo.IamUser, provisionInfo.Region, provisionInfo.BucketName, true)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(bucket); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func ProvisionAtlasClusterPost(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	provisionInfo := &provisioner.AtlasClusterProvisionInfo{}
	err := decoder.Decode(provisionInfo)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	client, err := provisioner.AtlasClient(viper.GetString("ATLAS_USERNAME"), viper.GetString("ATLAS_API_KEY"))
	atlastCluster, err := provisioner.NewCluster(client, viper.GetString("ATLAS_GROUP_ID"), provisionInfo)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(atlastCluster); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func ProvisionAtlasGetCluster(w http.ResponseWriter, r *http.Request) {

	clusterName := mux.Vars(r)["clusterName"]
	client, err := provisioner.AtlasClient(viper.GetString("ATLAS_USERNAME"), viper.GetString("ATLAS_API_KEY"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	mongoCluster, err := provisioner.GetCluster(client, viper.GetString("ATLAS_GROUP_ID"), clusterName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(mongoCluster); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

func ProvisionAtlasMongoUserPost(w http.ResponseWriter, r *http.Request) {

	userInfo := &provisioner.MongoDbProvisionInfo{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(userInfo)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Header().Add("Content-Type", "application/json")
	client, err := provisioner.AtlasClient(viper.GetString("ATLAS_USERNAME"), viper.GetString("ATLAS_API_KEY"))
	mongoDbUser, err := provisioner.NewMongoDbUser(client, viper.GetString("ATLAS_GROUP_ID"), userInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := json.NewEncoder(w).Encode(mongoDbUser); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

func ProvisionAtlasGetMongoDbUser(w http.ResponseWriter, r *http.Request) {

	userName := mux.Vars(r)["userName"]
	client, err := provisioner.AtlasClient(viper.GetString("ATLAS_USERNAME"), viper.GetString("ATLAS_API_KEY"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	mongoDbUser, err := provisioner.GetMongoDbUser(client, viper.GetString("ATLAS_GROUP_ID"), userName)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(mongoDbUser); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

func ProvisionersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
