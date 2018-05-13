package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/tomoya-k31/aws-sts-auth/sts"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	profile   string
	region    string
	tokenCode string

	awsIni         *sts.AwsIni
	awsIniSettings *sts.AwsIniSettings
)

func init() {
	// --profile default
	flag.StringVar(&profile, "profile", "default", "Named profile")
	// --region us-east-1
	flag.StringVar(&region, "region", "", "")
	// --token-code 123456
	flag.StringVar(&tokenCode, "token-code", "", "The value provided by the MFA device")
	flag.Parse()

	awsIni = sts.NewAwsIni()
	awsIniSettings = awsIni.Load(profile)
	if region == "" {
		if awsIniSettings.Region != "" {
			region = awsIniSettings.Region
		} else {
			region = "us-east-1"
		}
	}
}

func main() {

	var mfaSerialNumber string
	var mfaTokenCode string
	if terminal.IsTerminal(syscall.Stdin) {

		if awsIniSettings.MfaSerial == "" {
			fmt.Print("Type 'MFA serial number': ")
			fmt.Scan(&mfaSerialNumber)
			awsIniSettings.MfaSerial = mfaSerialNumber
		} else {
			fmt.Printf("Use [profile: %s] \"%s\"\n", profile, awsIniSettings.MfaSerial)
			mfaSerialNumber = awsIniSettings.MfaSerial
		}

		if tokenCode == "" {
			fmt.Print("Type \"MFA token code\": ")
			fmt.Scan(&mfaTokenCode)
		} else {
			mfaTokenCode = tokenCode
		}

	} else {
		fmt.Println("Fail to load pipe")
		os.Exit(1)
		return
	}

	token, err := sts.StsAuth(&sts.StsConfig{
		Region:          region,
		Profile:         profile,
		DurationSeconds: 3600, // 1h
		SerialNumber:    mfaSerialNumber,
		TokenCode:       mfaTokenCode,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	token.SaveEnv()
	token.SaveCredentials(profile, awsIni)
	token.Print()
}
