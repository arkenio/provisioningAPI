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

{ "name": "testnxio000003"}


Returns:

{"name":"testnxio000003","path":"","region":"us-east-1","iamUser":{"userName":"testnxio000003","userId":"AIDAJTSS7TGXVWD46JB7S","publicAccessKeyID":"AKIAIWKY24PN4UYFI24Q","publicAccessKey":"SHaZYOiSyieYcaWEsnOmwkKivcvX3ioiAX+FPZXl","arn":"arn:aws:iam::188670881089:user/testnxio000003"},"policyArn":"arn:aws:iam::188670881089:policy/nxios3policy"}{
   "name":"testnxio000003",
   "path":"",
   "region":"us-east-1",
   "iamUser":{
      "userName":"testnxio000003",
      "userId":"AIDAJTSS7TGXVWD46JB7S",
      "publicAccessKeyID":"AKIAIWKY24PN4UYFI24Q",
      "publicAccessKey":"SHaZYOiSyieYcaWEsnOmwkKivcvX3ioiAX+FPZXl",
      "arn":"arn:aws:iam::188670881089:user/testnxio000003"
   },
   "policyArn":"arn:aws:iam::188670881089:policy/nxios3policy"
}

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
