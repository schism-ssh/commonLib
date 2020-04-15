package commonLib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

func SSMClient(region string) ssmiface.SSMAPI {
	return ssm.New(AwsSession(region))
}

func LambdaClient(region string) lambdaiface.LambdaAPI {
	return lambda.New(AwsSession(region))
}

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
