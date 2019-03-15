# userpool-federation

Adds the ability to attach UserPool Federation through CloudFormation.

## Resource

The name for this custom resource is `Custom::CognitoUserPoolFederation` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::CognitoUserPoolFederation"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

These are the supported properties for the resource.

| Propertie name | Type | Description | Required |
| - | - | - | - |
| ProviderName | String | Name of the identity provider | Yes |
| UserPoolId | String | The ID of the UserPool to create the Identity Provider in | Yes |
| ProviderType | String | The Identity Provider Type. Valid options are: **SAML**, **Facebook**, **Google**, **LoginWithAmazon**, **OIDC** | Yes |
| ProviderDetails | List of strings | Details regarding your provider such as **MetadataURL**, **MetadataFile** etc. | Yes |
| AttributeMapping | List of strings | Identity Provider attribute mappings | No |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

For more details about the properties check the aws cli docs [https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-identity-provider.html](https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-identity-provider.html).

### ProviderName

**String** *Required*

Name of the Identity Provider you want to create.

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

  UserPoolFederationADFS:
    Type: "Custom::CognitoUserPoolFederation"
    DependsOn:
      - "UserPool"
    Properties:
      ProviderName: "ADFS"
      ProviderType: "SAML"
      ProviderDetails:
        MetadataURL: "https://my.domain.com/FederationMetadata.xml"
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-federation-${AWS::Region}-${Environment}"
      UserPoolId: !Ref "UserPool"
```