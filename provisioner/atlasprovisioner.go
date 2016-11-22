package provisioner

import (
	//"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/heimweh/go-mongodb/atlas"
	"time"
)

type AtlasProviderSettingsInfo struct {
	InstanceSizeName string `json:"instanceSizeName"`
	ProviderName     string `json:"providerName"`
	RegionName       string `json:"regionName"`
}

type AtlasClusterProvisionInfo struct {
	Name             string                     `json:"name"`
	BackupEnabled    bool                       `json:"backupEnabled"`
	ProviderSettings *AtlasProviderSettingsInfo `json:"providerSettings"`
}

type AtlasMongoDBCluster struct {
	Name              string                     `json:"name"`
	GroupId           string                     `json:"groupId"`
	MongoDBVersion    string                     `json:"mongoDBVersion"`
	MongoURI          string                     `json:"mongoURI"`
	MongoURIUpdated   time.Time                  `json:"mongoURIUpdated"`
	NumShards         int32                      `json:"numShards"`
	ReplicationFactor int32                      `json:"replicationFactor"`
	DiskSizeGB        int32                      `json:"diskSizeGB"`
	BackupEnabled     bool                       `json:"backupEnabled"`
	StateName         string                     `json:"stateName"`
	ProviderSettings  *AtlasProviderSettingsInfo `json:"providerSettings"`
}

type MongoDbProvisionInfo struct {
	DatabaseName string            `json:"databaseName"`
	Roles        []MongoDbUserRole `json:"roles"`
	Username     string            `json:"username"`
	Password     string            `json:"password"`
}

type MongoDbUserRole struct {
	DatabaseName string `json:"databaseName"`
	RoleName     string `json:"roleName"`
}

type MongoDBUser struct {
	Username     string            `json:"username"`
	GroupId      string            `json:"groupId"`
	Roles        []MongoDbUserRole `json:"roles"`
	DatabaseName string            `json:"databaseName"`
}

func AtlasClient(username string, apiKey string) (*atlas.Client, error) {
	client := atlas.NewClient(username, apiKey)
	glog.Infof("[INFO] MongoDB Atlas client configured")
	return client, nil
}

func GetGroupWhiteList(client *atlas.Client, groupID string) ([]atlas.GroupWhiteList, error) {
	list, _, err := client.GroupWhiteList.List(groupID)
	if err != nil {
		for _, wl := range list {
			glog.Infof("GroupID: %s", wl.GroupID)
			glog.Infof("CidrBlock: %s", wl.CidrBlock)
			glog.Infof("IPAddress: %s", wl.IPAddress)
		}
		return list, nil
	}
	return nil, err
}

//POST /api/atlas/v1.0/groups/GROUP-ID/clusters
//returns 200 if the cluster is provisioning
func NewCluster(client *atlas.Client, groupID string, newClusterInfo *AtlasClusterProvisionInfo) (*AtlasMongoDBCluster, error) {
	path := fmt.Sprintf("groups/%s/clusters", groupID)
	reqBody := newClusterInfo
	req, err := client.NewRequest("POST", path, reqBody)
	if err != nil {
		return nil, err
	}
	glog.Errorf("Req body %s", reqBody)
	mogoCluster := new(AtlasMongoDBCluster)
	_, err = client.Do(req, mogoCluster)

	if err != nil {
		glog.Infof("Error while trying to provision server %v", err)
		return nil, err
	}
	/**
	Sample response from the server:
	{
	 "backupEnabled":true,
	 "diskSizeGB":40,
	 "groupId":"57d8232ed383ad6a442810bb",
	 "links":[{"href":"https://cloud.mongodb.com/api/atlas/v1.0/groups/57d8232ed383ad6a442810bb/clusters/testprovisioning5","rel":"self"}],
	 "mongoDBMajorVersion":"3.2",
	 "mongoDBVersion":"3.2.10",
	 "mongoURIUpdated":"2016-11-13T23:26:06Z"
	 ,"name":"testprovisioning5",
	 "numShards":1,
	 "providerSettings":
	     {
	      "providerName":"AWS",
	      "diskIOPS":120,
	      "encryptEBSVolume":false,
	      "instanceSizeName":"M10",
	      "regionName":"US_EAST_1"
	      },
	 "replicationFactor":3,
	 "stateName":"CREATING"}
	**/

	return mogoCluster, nil
}

func GetCluster(client *atlas.Client, groupID string, clusterName string) (*AtlasMongoDBCluster, error) {
	path := fmt.Sprintf("groups/%s/clusters/%s", groupID, clusterName)
	req, err := client.NewRequest("GET", path, "")
	if err != nil {
		return nil, err
	}
	mogoCluster := new(AtlasMongoDBCluster)
	_, err = client.Do(req, mogoCluster)

	if err != nil {
		glog.Infof("Error while tryin to fetch cluster %v %s", err)
		return nil, err
	}
	return mogoCluster, nil
}

func NewMongoDbUser(client *atlas.Client, groupID string, newMongoDbUserInfo *MongoDbProvisionInfo) (*MongoDBUser, error) {

	path := fmt.Sprintf("groups/%s/databaseUsers/", groupID)
	roles := make([]MongoDbUserRole, 0)
	roles = append(roles, MongoDbUserRole{newMongoDbUserInfo.DatabaseName, "readWrite"})
	newMongoDbUserInfo.DatabaseName = "admin"
	newMongoDbUserInfo.Roles = roles

	reqBody := newMongoDbUserInfo

	req, err := client.NewRequest("POST", path, reqBody)
	if err != nil {
		return nil, err
	}
	mogoDbUser := new(MongoDBUser)
	_, err = client.Do(req, mogoDbUser)

	if err != nil {
		glog.Infof("Error while trying to provision user: %v", err)
		return nil, err
	}
	return mogoDbUser, nil
}

func GetMongoDbUser(client *atlas.Client, groupID string, username string) (*MongoDBUser, error) {

	path := fmt.Sprintf("groups/%s/databaseUsers/admin/%s", groupID, username)
	req, err := client.NewRequest("GET", path, "")
	if err != nil {
		return nil, err
	}
	mogoDbUser := new(MongoDBUser)
	_, err = client.Do(req, mogoDbUser)

	if err != nil {
		glog.Infof("Error while trying to fetch user: %v", err)
		return nil, err
	}
	return mogoDbUser, nil
}
