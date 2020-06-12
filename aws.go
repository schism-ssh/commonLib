package commonLib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// SSMClient returns a new AWS SSM Client in a given region
func SSMClient(region string) ssmiface.SSMAPI {
	return ssm.New(AwsSession(region))
}

// LambdaClient returns a new AWS Lambda Client in a given region
func LambdaClient(region string) lambdaiface.LambdaAPI {
	return lambda.New(AwsSession(region))
}

// S3Client returns a new AWS S3 Client in a given region
func S3Client(region string) s3iface.S3API {
	return s3.New(AwsSession(region))
}

// AwsSession returns an active session for AWS APIs given a region.
//
// This also hard-enables SharedConfig access for Schism AWS API access.
//
// TODO: Make this configurable somehow. (Lambda vs CLI Tool?)
func AwsSession(region string) *session.Session {
	sessionOpts := session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigEnable,
	}
	awsSession := session.Must(session.NewSessionWithOptions(sessionOpts))
	return awsSession
}
