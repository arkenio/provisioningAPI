### Running the server
To run the server, follow these simple steps:

```
go run main.go
```

### Provision IAM user and S3 bucket
Provision the specified IAM user and S3 bucket with the given name in the given region.
Configure specific bucket policy.
Expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env variables to be set for now.



POST http://localhost:8080/provision/s3
body:

```{ 
    "bucketName": "testnxio000004",
    "iamUser": "testnxio000004User",
    "region": "us-east-1"
    }
```

If these are not specified, fallback on default values

Returns:
```{
  "name": "testnxio000004",
  "path": "",
  "region": "us-east-1",
  "iamUser": {
    "userName": "testnxio000004User",
    "userId": "AIDAIEAAUSHYTZOZGOUXG",
    "publicAccessKeyID": "***",
    "publicAccessKey": "***",
    "arn": "arn:aws:iam::188670881089:user/testnxio000004User"
  },
  "policyArn": "arn:aws:iam::188670881089:policy/nxios3policy"
}
```

### Provision an AtlasMongoDB cluster and a dedicated MongoDB user for the database
Expects ATLAS_USERNAME, ATLAS_API_KEY and ATLAS_GROUP_ID to set as env variables

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
