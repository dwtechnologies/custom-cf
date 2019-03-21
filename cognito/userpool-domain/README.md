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

| Property name | Type | Description | Required |
| - | - | - | - |
| CSS | String | CSS to use for the UI | No |
| ImageFile | String | Base64 encoded Image | Yes if custom domain is used |
| ClientId | String | The UserPool Client ID | Yes |
| UserPoolId | String | The ID of the UserPool to create the Identity Provider in | Yes |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

For more details about the custom domain check [https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/set-ui-customization.html](https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/set-ui-customization.html).

## Supported Attributes

The following attributes can be used in CloudFormations `Fn::GetAtt` function.

- CSSVersion
- ClientId
- UserPoolId

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

  UserPoolDomain:
    Type: "Custom::CognitoUserPoolDomain"
    DependsOn:
      - "UserPool"
    Properties:
      Domain: "mydevtestpoolcustom"
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-domain-${AWS::Region}-${Environment}"
      UserPoolId: !Ref "UserPool"
```