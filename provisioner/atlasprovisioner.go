package provisioner

import (
	//"encoding/json"
	"fmt"
	"github.com/arkenio/provisioningAPI/go-mongodb/atlas"
	"github.com/golang/glog"
	"time"
)

type AtlasProviderSettingsInfo struct {
	InstanceSizeName string `json:"instanceSizeName,omitempty"`
	ProviderName     string `json:"providerName"`
	RegionName       string `json:"regionName,omitempty"`
	DiskIOPS         int32 `json:"diskIOPS,omitempty"`
	EncryptEBSVolume bool `json:"encryptEBSVolume,omitempty"`
}

type AtlasClusterProvisionInfo struct {
	Name             string                     `json:"name,omitempty"`
	BackupEnabled    bool                       `json:"backupEnabled"`
	ProviderSettings *AtlasProviderSettingsInfo `json:"providerSettings"`
}

type AtlasMongoDBCluster struct {
	Name              string                     `json:"name,omitempty"`
	GroupId           string                     `json:"groupId"`
	MongoDBVersion    string                     `json:"mongoDBVersion,omitempty"`
	MongoURI          string                     `json:"mongoURI,omitempty"`
	MongoURIUpdated   time.Time                  `json:"mongoURIUpdated,omitempty"`
	NumShards         int32                      `json:"numShards,omitempty"`
	ReplicationFactor int32                      `json:"replicationFactor,omitempty"`
	DiskSizeGB        int32                      `json:"diskSizeGB,omitempty"`
	BackupEnabled     bool                       `json:"backupEnabled,omitempty"`
	StateName         string                     `json:"stateName,omitempty"`
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

func ModifyCluster(client *atlas.Client, groupID string, clusterName string, newClusterInfo *AtlasClusterProvisionInfo) (*AtlasMongoDBCluster, error) {
	path := fmt.Sprintf("groups/%s/clusters/%s", groupID, clusterName)
	reqBody := newClusterInfo
	glog.Errorf("Req body %s", reqBody)
	req, err := client.NewRequest("PATCH", path, reqBody)
	if err != nil {
		return nil, err
	}
	mogoCluster := new(AtlasMongoDBCluster)
	_, err = client.Do(req, mogoCluster)

	if err != nil {
		glog.Infof("Error while trying to provision server %v", err)
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
