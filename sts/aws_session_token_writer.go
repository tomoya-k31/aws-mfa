package sts

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ini/ini"
)

func (token *AwsSessionToken) SaveEnv() {
	os.Setenv("AWS_ACCESS_KEY_ID", token.AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", token.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", token.SessionToken)
}

func (token *AwsSessionToken) SaveCredentials(profile string, awsIni *AwsIni) {
	filepath, _ := awsIni.getFilepath(Credentials)
	cfg, err := ini.Load(filepath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	tempProfile := profile + "-sts"
	cfg.Section(tempProfile).Key("aws_access_key_id").SetValue(token.AccessKeyId)
	cfg.Section(tempProfile).Key("aws_secret_access_key").SetValue(token.SecretAccessKey)
	cfg.Section(tempProfile).Key("aws_session_token").SetValue(token.SessionToken)

	if err := cfg.SaveTo(filepath); err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	os.Setenv("AWS_PROFILE", tempProfile)
}

func (token *AwsSessionToken) Print() {
	fmt.Printf("Credentials[AccessKeyId]:     \"%s\"\n", token.AccessKeyId)
	fmt.Printf("Credentials[SecretAccessKey]: \"%s\"\n", token.SecretAccessKey)
	fmt.Printf("Credentials[SessionToken]:    \"%s\"\n", token.SessionToken)
	fmt.Printf("Credentials[Expiration]:      \"%s\"\n", token.Expiration.Format(time.RFC3339))
}
