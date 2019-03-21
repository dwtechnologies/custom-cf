# userpool-mfa

Adds the ability to create, update and delete UserPool MFA Settings through CloudFormation.

## Resource

The name for this custom resource is `Custom::CognitoUserPoolMFA` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::CognitoUserPoolMFA"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

| Property name | Type | Description | Required |
| - | - | - | - |
| MfaConfiguration | String | If MFA should be enabled. Possible values OFF, ON, OPTIONAL | Yes |
| UserPoolId | String | The ID of the UserPool to create the Identity Provider in | Yes |
| SmsMfaConfiguration | SmsMfaConfiguration | The SMS configuration if MFA should be via SMS | No |
| SoftwareTokenMfaConfiguration | SoftwareTokenMfaConfiguration | The Software Token configuration if MFA should be via software | No |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

For more details about the properties check the aws cli docs [https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/set-user-pool-mfa-config.html](https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/set-user-pool-mfa-config.html).

### SmsMfaConfiguration Properties

| Property name | Type | Description | Required |
| - | - | - | - |
| SmsAuthenticationMessage | String | SMS message to send for authentication | Yes |
| SmsConfiguration | SmsConfiguration | Configuration for sending SMS through AWS | Yes |

### SmsConfiguration Properties

| Property name | Type | Description | Required |
| - | - | - | - |
| SnsCallerArn | String | ARN to the SNS caller | Yes |
| ExternalId | String | The external ID | No |

### SoftwareTokenMfaConfiguration Properties

| Property name | Type | Description | Required |
| - | - | - | - |
| Enabled | bool | If Software OTP should be enabled | Yes |

## Example

```yaml
AWSTemplateFormatVersion: "2010-09-09"
Description: "Cognito UserPool"

Parameters:
  Environment:
    Description: "What environment we deploy to"
    Type: "String"
    Default: "dev"

Resources:
  UserPool:
    Type: "AWS::Cognito::UserPool"
    Properties:
      AliasAttributes:
        - "email"
      MfaConfiguration: "OFF"
      UserPoolName: "userpool"

  UserPoolClient:
    Type: "Custom::CognitoUserPoolMFA"
    DependsOn:
      - "UserPool"
    Properties:
      MfaConfiguration: "ON"
      SoftwareTokenMfaConfiguration:
        Enabled: true
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-mfa-${AWS::Region}-${Environment}"
      UserPoolId: !Ref "UserPool"
```