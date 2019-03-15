# userpool-domain

Adds the ability to create a UserPool Domain through CloudFormation.

## Resource

The name for this custom resource is `Custom::CognitoUserPoolDomain` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::CognitoUserPoolDomain"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

These are the supported properties for the resource.

| Propertie name | Type | Description | Required |
| - | - | - | - |
| Domain | String | Domain name or part of the domain name to use | Yes |
| UserPoolId | String | The ID of the UserPool to create the Identity Provider in | Yes |
| CustomDomainConfig | Object | Object with CustomDomainConfig | Yes if custom domain is used |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

For more details about the custom domain check [https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-assign-domain.html](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-assign-domain.html).

### CustomDomainConfig Properties

| Propertie name | Type | Description | Required |
| - | - | - | - |
| CertificateArn | String | ARN to ACM Certificate | Yes |

## Example

```yaml
AWSTemplateFormatVersion: "2010-09-09"
Description: "Cognito UserPool with SAML federation"

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

  UserPoolDomain:
    Type: "Custom::CognitoUserPoolDomain"
    DependsOn:
      - "UserPool"
    Properties:
      Domain: "mydevtestpoolcustomdw"
      ProviderType: "SAML"
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-domain-${AWS::Region}-${Environment}"
      UserPoolId: !Ref "UserPool"
```