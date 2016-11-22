package provisioner

import (
	aws "github.com/aws/aws-sdk-go/aws"
	awserr "github.com/aws/aws-sdk-go/aws/awserr"
	cred "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
	"strings"
)

type S3ProvisionInfo struct {
	BucketName string `json:"bucketName"`
	IamUser    string `json:"iamUser"`
	Region     string `json:"region"`
}

type IAMUser struct {
	UserName          string `json:"userName"`
	UserId            string `json:"userId"`
	PublicAccessKeyID string `json:"publicAccessKeyID"`
	PublicAccessKey   string `json:"publicAccessKey"`
	Arn               string `json:"arn"`
}

type S3Bucket struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Region    string `json:"region,omitempty""`
	IamUser   *IAMUser `json:"iamUser""`
	PolicyArn string `json:"policyArn""`
}

const (
	DEFAULT_S3IO_BUCKET_NAME             = "testnxios3bucket"
	DEFAULT_S3IO_IAMUSERNAME             = "testnxios3bucket_iam"
	DEFAULT_S3IO_ACCESSBUCKET_POLICYNAME = "testnxios3policy"
	DEFAULT_AWS_S3_REGION                = "us-east-1"
	DEFAULT_AWS_S3_POLICY_TEMPLATE       = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListAllMyBuckets"
            ],
            "Resource": "arn:aws:s3:::*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetBucketLocation",
                "s3:AbortMultipartUpload",
                "s3:ListMultipartUploadParts",
                "s3:ListBucketMultipartUploads"
            ],
            "Resource": "arn:aws:s3:::{bucket}/{path}"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject",
                "s3:DeleteObject",
                "s3:AbortMultipartUpload",
                "s3:ListMultipartUploadParts",
                "s3:ListBucketMultipartUploads"
            ],
            "Resource": "arn:aws:s3:::{bucket}/{path}/*"
        }
    ]
}`
)

func ProvisionIAMUserIfDoesntExist(session *awssession.Session, userName string, configNewAccessKey bool) (*IAMUser, error) {
	iamUser := &IAMUser{}
	svcIam := iam.New(session)

	// provision a new IAM user for this, if doesn't exist

	userInfo, err := svcIam.GetUser(&iam.GetUserInput{
		UserName: aws.String(userName),
	})

	if err == nil {
		glog.Infof("User %s already exists, no need to provision", *userInfo.User.UserName)
		iamUser.UserName = *userInfo.User.UserName
		iamUser.UserId = *userInfo.User.UserId
		iamUser.Arn = *userInfo.User.Arn

	}

	if err != nil && err.(awserr.Error).Code() == "NoSuchEntity" {
		userInfo, err := svcIam.CreateUser(&iam.CreateUserInput{
			UserName: aws.String(userName),
		})

		if err != nil {
			glog.Errorf("Error creating user, %v", err)
			return nil, err
		}

		iamUser.UserName = *userInfo.User.UserName
		iamUser.UserId = *userInfo.User.UserId
		iamUser.Arn = *userInfo.User.Arn
		glog.Infof("Provisioned user %s", iamUser.UserName)
	}

	if configNewAccessKey {
		glog.Infof("Removing all the previous access keys for %s and generating a new one", iamUser.UserName)
		//if true deletes any existing access key if any for this user

		//list all the access keys
		keys, err := svcIam.ListAccessKeys(&iam.ListAccessKeysInput{
			UserName: aws.String(iamUser.UserName),
		})

		for _, key := range keys.AccessKeyMetadata {
			_, err = svcIam.DeleteAccessKey(&iam.DeleteAccessKeyInput{
				UserName:    aws.String(iamUser.UserName),
				AccessKeyId: key.AccessKeyId,
			})
		}

		// Give the user a new key
		res, err := svcIam.CreateAccessKey(&iam.CreateAccessKeyInput{
			UserName: aws.String(iamUser.UserName),
		})
		if err != nil {
			glog.Errorf("Can not generate access key %v", err)
			return iamUser, err
		}

		iamUser.PublicAccessKeyID = *res.AccessKey.AccessKeyId
		iamUser.PublicAccessKey = *res.AccessKey.SecretAccessKey

	}

	return iamUser, nil
}

func AttachS3BucketPolicyToIAMUser(session *awssession.Session, bucket *S3Bucket, iamUser *IAMUser) (*S3Bucket, error) {

	svcIam := iam.New(session)
	user_policy := strings.Replace(DEFAULT_AWS_S3_POLICY_TEMPLATE, "{bucket}", bucket.Name, -1)
	user_policy = strings.Replace(user_policy, "/{path}", bucket.Path, -1)

	res, err := svcIam.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(user_policy),
		PolicyName:     aws.String(DEFAULT_S3IO_ACCESSBUCKET_POLICYNAME),
	})
	if err == nil {
		bucket.PolicyArn = *res.Policy.Arn
	}

	if err != nil && err.(awserr.Error).Code() == "EntityAlreadyExists" {
		//policy alreday exists
		//it might be already attached to this user
		userpolicies, err := svcIam.ListPolicies(&iam.ListPoliciesInput{})
		if err != nil {
			return bucket, err
		}
		for _, policy := range userpolicies.Policies {

			if *policy.PolicyName == DEFAULT_S3IO_ACCESSBUCKET_POLICYNAME {
				bucket.PolicyArn = *policy.Arn
				break
			}
		}
	}
	//should not break even already attached
	_, err = svcIam.AttachUserPolicy(&iam.AttachUserPolicyInput{
		PolicyArn: aws.String(bucket.PolicyArn),
		UserName:  aws.String(iamUser.UserName),
	})

	if err != nil {
		glog.Errorf("Failed to attach policy: %s to user: %s, with error %v", bucket.PolicyArn, bucket.IamUser.UserName, err)
		return bucket, err
	}

	return bucket, nil
}

func ProvisionBucketIfDoesntExist(session *awssession.Session, bucketName string, bucketPath string, region string) (*S3Bucket, error) {

	bucket := &S3Bucket{}
	svc := s3.New(session)

	bkt, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: &bucketName,
	})

	if err != nil && err.(awserr.Error).Code() == "BucketAlreadyOwnedByYou" {
		glog.Infof("BucketAlreadyOwnedByYou %s, no need to reprovision", bucketName)
	}
	if err != nil && err.(awserr.Error).Code() != "BucketAlreadyOwnedByYou" {
		glog.Infof("Error while provisioning bucket %v", err)
		return nil, err
	}

	if err == nil && bkt != nil && bkt.Location != nil {
		glog.Infof("Provisioned bucket %s at location %s", bucketName, *bkt.Location)
		bucket.Name = bucketName
		bucket.Path = bucketPath
		bucket.Region = region
	}

	if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucketName}); err != nil {
		glog.Infof("Failed to wait for bucket to exist %s, %s\n", bucketName, err)
		return bucket, err
	}

	return bucket, nil
}

func (bucket *S3Bucket) ProvisionS3AndIAMUser(username string, region string, bucketname string, configNewAccessKey bool) (*S3Bucket, error) {
	awsS3Regions := []string{"us-east-1", "us-west-1", "us-west-2", "eu-west-1", "eu-central-1",
		"ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "sa-east-1"}

	validRegion := false
	for _, r := range awsS3Regions {
		if r == region {
			validRegion = true
			break
		}
	}

	if !validRegion {
		glog.Infof("%s is not a valid regsion, fallback on the default region", region)
		region = DEFAULT_AWS_S3_REGION
	}

	glog.Infof("Provisioning in %s", region)
	//expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env variables to be set
	session := awssession.New(&aws.Config{Region: aws.String(region), Credentials: cred.NewEnvCredentials()})
	iamUser, err := ProvisionIAMUserIfDoesntExist(session, username, configNewAccessKey)
	if err != nil {
		glog.Error("Failed to provision or fetch existing IAM user %v", err)
		return nil, err
	}
	glog.Infof("IAM user %v provisioned", iamUser)

	bucket, err = ProvisionBucketIfDoesntExist(session, bucketname, "", region)
	if err != nil {
		glog.Error("Failed to provision or fetch existing bucket %v", err)
		return nil, err
	}

	bucket.IamUser = iamUser
	_, err = AttachS3BucketPolicyToIAMUser(session, bucket, iamUser)
	if err != nil {
		glog.Error("Failed to attach security policy to iam user %v", err)
		return nil, err
	}
	return bucket, nil
}
