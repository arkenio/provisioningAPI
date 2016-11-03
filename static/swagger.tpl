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
       Expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set as env vars or
       to be sent in the body ( not cool)
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
      AWS_ACCESS_KEY_ID:
        type: string
      AWS_SECRET_ACCESS_KEY:
        type: string
      Name:
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
                  