# AWS MFA

## Overview
aws-mfa is a tool to create temporary session tokens using MFA devices.

## TODO
- [ ] test code
- [ ] CI(build status badge)
- [ ] dep libs
- [ ] distribution
- [ ] Makefile

## Installing

```shell
go get -u github.com/tomoya-k31/aws-mfa
```

## Using

```shell
$ aws-mfa --profile ex_user
Use [profile: ex_user] "arn:aws:iam::xxxxxxxxxxxx:mfa/ex_user"
Type "MFA token code": 999999    # <- type!!
Credentials[AccessKeyId]:     "....."
Credentials[SecretAccessKey]: ".........."
Credentials[SessionToken]:    "....................."
Credentials[Expiration]:      "2018-06-01T09:27:20Z"
```

or 

```shell
$ aws-mfa --profile ex_user --token-code 999999
```


## Creds/Config file format

ex) profile = ex_user

- `~/.aws/config`

```ini
[profile ex_user]
region = ap-northeast-1
mfa_serial = arn:aws:iam::xxxxxxxxxxxx:mfa/ex_user
```


- `~/.aws/credentials`

```ini
[ex_user]
aws_access_key_id     = XXXXXXXXXXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

## Ref

- [AWS CLI Configuration Variables — AWS CLI 1.15.30 Command Reference](https://docs.aws.amazon.com/cli/latest/topic/config-vars.html#credentials)
