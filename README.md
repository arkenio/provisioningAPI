### Running the server
To run the server, follow these simple steps:

Run:
```
go run main.go
```
Or build Docker image:
```
docker build -t arkenio/provisioner:v1
```
Or run Docker container:
```
docker run arkenio/provisioner:v1
```

### Provision a S3 bucket and a dedicated IAM user
Provision the specified IAM user and S3 bucket with the given name in the given region (this will also configure a bucket policy for the user to access the bucket as a inner policy to the user).

Expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env variables to be set for now.



POST http://localhost:8788/provision/s3
body:

```
{ 
    "bucketName": "testnxio000004",
    "iamUser": "testnxio000004User",
    "region": "us-east-1"
}
```

If these are not specified, fallback on default values.

Returns:
```
{
  "name": "testnxio000004",
  "path": "",
  "region": "us-east-1",
  "iamUser": {
    "userName": "testnxio000004User",
    "userId": "AIDAIEAAUSHYTZOZGOUXG",
    "publicAccessKeyID": "***",
    "publicAccessKey": "***",
    "arn": "arn:aws:iam::188670881089:user/testnxio000004User"
  }
}
```

To configure Nuxeo to use your provisioned resources set:
       nuxeo.s3storage.bucket=$name, 
       nuxeo.s3storage.awsid=$publicAccessKeyID,
       nuxeo.s3storage.awssecret=$publicAccessKey


### Provision an AtlasMongoDB cluster and a dedicated MongoDB user for the database
Expects ATLAS_USERNAME, ATLAS_API_KEY and ATLAS_GROUP_ID to set as env variables
You need to provision a cluster and a dedicated mongoDb user for your database. ( Nuxeo will automatically create the database and the collections at start-up)

        
1. Provision a cluster

POST http://localhost:8788/provision/atlas/clusters
body:

```
{
  "Name": "testprovisioning1",
  "BackupEnabled": true,
  "ProviderSettings": {
    "InstanceSizeName": "M10",
    "ProviderName": "AWS",
    "RegionName": "US_EAST_1"
  }
}
```
Returns:
```
{
  "name": "testprovisioning1",
  "groupId": "57d8232ed383ad6a442810bb",
  "mongoDBVersion": "3.2.11",
  "mongoURI": "",
  "numShards": 1,
  "replicationFactor": 3,
  "diskSizeGB": 40,
  "backupEnabled": true,
  "stateName": "CREATING",
  "providerSettings": {
    "instanceSizeName": "M10",
    "providerName": "AWS",
    "regionName": "US_EAST_1"
  }
}
```
        
 2. Provision a MongoDb user:
 
 POST http://localhost:8788/atlas/users
 
 ```
 {
  "username": "testprovisioning1",
  "databaseName": "nuxeo",
  "password": "nuxeo"
 }
```
Returns:

```
{
  "username": "testprovisioning1",
  "groupId": "57d8232ed383ad6a442810bb",
  "roles": [
    {
      "databaseName": "nuxeo",
      "roleName": "readWrite"
    }
  ],
  "databaseName": "admin"
}
```
3. Fetch back the cluster. If the instance has been created, the mongoURI will be returned
       
GET http://localhost:8788/provision/atlas/clusters/{testprovisioning1}

returns also:
```
{
"mongoURI": "mongodb://testprovisioning1-shard-00-00-buc6y.mongodb.net:27017,testprovisioning1-shard-00-01-buc6y.mongodb.net:27017,testprovisioning1-shard-00-02-buc6y.mongodb.net:27017",
}
```

For the above example, configure Nuxeo as:
  nuxeo.mongodb.dbname=nuxeo  
  nuxeo.mongodb.server= mongodb://nuxeo:nuxeo@testprovisioning1-shard-00-00-buc6y.mongodb.net:27017,testprovisioning1-shard-00-01-buc6y.mongodb.net:27017,testprovisioning1-shard-00-02-buc6y.mongodb.net:27017
* generally the nuxeo.mongodb.server = sprintf("mongodb://%s:%s@%s", $mongoDBUser, $mongoDBDatabase, $mongoURI)



About Nuxeo
-----------

Nuxeo provides a modular, extensible Java-based
[open source software platform for enterprise content management](http://www.nuxeo.com/en/products/ep),
and packaged applications for [document management](http://www.nuxeo.com/en/products/document-management),
[digital asset management](http://www.nuxeo.com/en/products/dam) and
[case management](http://www.nuxeo.com/en/products/case-management).

Designed by developers for developers, the Nuxeo platform offers a modern
architecture, a powerful plug-in model and extensive packaging
capabilities for building content applications.

More information on: <http://www.nuxeo.com/>
