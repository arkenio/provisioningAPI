# Example YAML to get you started quickly.
# Be aware that YAML has indentation based scoping.
# Code completion support is available so start typing for available options.
swagger: '2.0'

# This is your document metadata
info:
  version: "0.0.0"
  title: Provisioning service

schemes:
  - http

# Describe your paths here
paths:
  /provision/s3:
   post:
     description :
       This provision a new s3 bucket and a new IAM user with the associated policy.
       Expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set as env vars( not cool)
     parameters:
        - in: "body"
          name: "body"
          description: "Service definition that has to be created"
          required: true
          schema:
            $ref: "#/definitions/S3BucketForCreation"
     responses:
       200:
         description: Sucessfuly provisioned
         schema:
            $ref: "#/definitions/S3Bucket"
  /provision/atlas/clusters:
   post:
     description: 
        Provsion a new Atlas MongoDB Cluster.
        Expects ATLAS_USERNAME, ATLAS_GROUP_ID and ATLAS_API_KEY to be set as env vars
     parameters:
        - in: "body"
          name: "body"
          description: "Service definition that has to be created"
          required: true
          schema:
            $ref: "#/definitions/AtlasClusterForCreation"
     responses:
       200:
         description: Sucessfuly provisioned
         schema:
            $ref: "#/definitions/AtlasMongoDbCluster"
  /provision/atlas/clusters/{clusterName}:
   get:
    description: |
        Gets an existing cluster.
         Expects ATLAS_USERNAME, ATLAS_GROUP_ID and ATLAS_API_KEY to be set as env vars.
    parameters:
        - name: clusterName
          in: path
          required: true
          type: string
    responses:
       200:
         description: Returns the cluster
         schema:
            $ref: "#/definitions/AtlasMongoDbCluster"            
  /provision/atlas/users:
   post:
     description: 
        Provsion a new Atlas MongoDB User for any cluster in the group for the given database
     parameters:
        - in: "body"
          name: "body"
          description: "The name/password of the user and the name of the database"
          required: true
          schema:
            $ref: "#/definitions/MongoDbUserForCreation"
     responses:
       200:
         description: Sucessfuly provisioned
         schema:
            $ref: "#/definitions/MongoDBUser"
  /provision/atlas/users/{userName}:
   get:
    description: |
        Gets an existing MongoDb user.
         Expects ATLAS_USERNAME, ATLAS_GROUP_ID and ATLAS_API_KEY to be set as env vars.
    parameters:
        - name: userName
          in: path
          required: true
          type: string
    responses:
       200:
         description: Returns the mongoDb user
         schema:
            $ref: "#/definitions/MongoDBUser"           
       
  # This is a path endpoint. Change it.
  /provisioners:
    # This is a HTTP operation
    get:
      # Describe this verb here. Note: you can use markdown
      description: |
        This could get the list of provisioners. Not implemented for now
    
      # Expected responses for this operation:
      responses:
        # Response code
        200:
          description: Successful response
          # A schema describing your response object.
          # Use JSON Schema format
          schema:
            title: ArrayOfProvisioners
            type: array
            items:
              title: Provisioner
              type: object
              properties:
                name:
                  type: string
                single:
                  type: boolean
definitions:
  S3BucketForCreation:
    type: object
    properties:
      BucketName:
        type: string
      IamUser:
        type: string
      Region:
         type: string
  S3Bucket:
    type: object
    properties:
      name:
        type: string
      path:
        type: string
      region:
         type: string
      iamUser:
        $ref: '#/definitions/IAMUser'
  IAMUser:
    type: object
    properties:
      userName:
        type: string
      UserId:
        type: string
      PublicAccessKeyID:
        type: string
      PublicAccessKey:
         type: string 
      Arn:
         type: string
  AtlasClusterForCreation:
    type: object
    properties:
      Name:
        type: string
      BackupEnabled:
        type: boolean
      ProviderSettings:
        $ref: '#/definitions/AtlastProviderSettingsForCreation'
    example:
        Name: testprovisioning1
        BackupEnabled: true
        ProviderSettings:
            InstanceSizeName: M10
            ProviderName: AWS
            RegionName: US_EAST_1
  AtlastProviderSettingsForCreation:
    type: object
    properties:
      InstanceSizeName:
        type: string
      ProviderName:
        type: string
      RegionName:
        type: string
    example:
      InstanceSizeName: M10
      ProviderName: AWS
      RegionName: US_EAST_1
  AtlastProviderSettings:
    type: object
    properties:
      instanceSizeName:
        type: string
      providerName:
        type: string
      regionName:
        type: string
      diskIOPS:
         type: integer
      encryptEBSVolume:
         type: boolean
  AtlasMongoDbCluster:
    type: object
    properties:
      name:
        type: string
      groupId:
        type: string
      mongoDBVersion:
        type: string
      mongoURI:
         type: string 
      mongoURIUpdated:
         type: string
      numShards:
         type: integer
      replicationFactor:
         type: integer
      providerSettings:   
         $ref: '#/definitions/AtlastProviderSettings'
      diskSizeGB:
          type: integer
      backupEnabled:
          type: boolean
      stateName:
          type: string
  MongoDbUserForCreation:
    type: object
    properties:
      username:
        type: string
      databaseName:
        type: string
      password:
        type: string  
    example:
      username: testprovisioning
      databaseName: nuxeo
      password: nuxeo
  MongoDBUser:
    type: object
    properties:
      databaseName:
        type: string
      groupId:
        type: string
      username:
        type: string
      roles:
        type: array
        items:
          type: object
          properties:
           databaseName:
             type: string
           roleName:
              type: string