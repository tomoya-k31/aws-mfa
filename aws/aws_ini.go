package aws

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	"github.com/go-ini/ini"
)

// FileType Enum
type FileType int8

const (
	Config FileType = iota
	Credentials
)

var (
	homeDir string
)

// AwsIni AWS設定ファイル
type AwsIni struct {
	Config      string
	Credentials string
	AwsDir      string
}

type AwsIniSettings struct {
	Region             string
	MfaSerial          string
	AwsAccessKeyId     string
	AwsSecretAccessLey string

	// ConfigFilename      string
	// CredentialsFilename string
}

func init() {
	// $HOME
	if user, err := user.Current(); err != nil {
		log.Fatal(err)
	} else {
		homeDir = user.HomeDir
	}
}

func NewAwsIni() *AwsIni {
	return &AwsIni{
		Config:      "config",
		Credentials: "credentials",
		AwsDir:      ".aws",
	}
}

func (awsIni *AwsIni) Load(profile string) *AwsIniSettings {
	settings := &AwsIniSettings{}

	if section, err := awsIni.loadIniSection(profile, Config); err == nil {
		settings.Region = section.Key("region").String()
		settings.MfaSerial = section.Key("mfa_serial").String()
	}

	if section, err := awsIni.loadIniSection(profile, Credentials); err == nil {
		settings.AwsAccessKeyId = section.Key("aws_access_key_id").String()
		settings.AwsSecretAccessLey = section.Key("aws_secret_access_key").String()
	}

	return settings
}

func (awsIni *AwsIni) loadIniSection(profile string, f FileType) (section *ini.Section, err error) {
	filepath, _ := awsIni.getFilepath(f)
	if cfg, err := ini.Load(filepath); err == nil {
		section = cfg.Section(getAwsIniSection(profile, f))
	}
	return
}

func (awsIni *AwsIni) getFilepath(f FileType) (s string, err error) {
	switch f {
	case Config:
		s = filepath.Join(homeDir, awsIni.AwsDir, awsIni.Config)
	case Credentials:
		s = filepath.Join(homeDir, awsIni.AwsDir, awsIni.Credentials)
	default:
		err = fmt.Errorf("Cannot create file path: unknown filetype")
	}
	return
}

// GetAwsIniSection get aws section name in config-ini
func getAwsIniSection(profile string, f FileType) string {
	switch {
	case f == Config && profile != "default":
		return "profile " + profile
	case f == Config && profile == "default":
		fallthrough
	case f == Credentials:
		fallthrough
	default:
		return profile
	}
}
