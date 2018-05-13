package awsmfa

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type StsConfig struct {
	Region          string
	Profile         string
	DurationSeconds int64
	SerialNumber    string
	TokenCode       string
}

type AwsSessionToken struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

func StsAuth(s *StsConfig) (*AwsSessionToken, error) {

	// creds := credentials.NewSharedCredentials("", profile)
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{
				Filename: "",
				Profile:  s.Profile,
			},
			&credentials.EnvProvider{},
		})

	config := aws.Config{Region: aws.String(s.Region), Credentials: creds}
	sess, err := session.NewSession(&config)
	if err != nil {
		return nil, err
	}

	svc := sts.New(sess)
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(s.DurationSeconds),
		SerialNumber:    aws.String(s.SerialNumber),
		TokenCode:       aws.String(s.TokenCode),
	}

	result, err := svc.GetSessionToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeRegionDisabledException:
				// fmt.Println(sts.ErrCodeRegionDisabledException, aerr.Error())
				return nil, fmt.Errorf(sts.ErrCodeRegionDisabledException + " " + aerr.Error())
			default:
				// fmt.Println(aerr.Error())
				return nil, aerr
			}
		} else {
			return nil, err
		}
	}

	return &AwsSessionToken{
		AccessKeyId:     aws.StringValue(result.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(result.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(result.Credentials.SessionToken),
		Expiration:      aws.TimeValue(result.Credentials.Expiration),
	}, nil
}
