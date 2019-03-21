# userpool-uicustomization

Adds the ability to set Custom UI elements to your UserPool through CloudFormation.

## Resource

The name for this custom resource is `Custom::CognitoUserPoolUICustomization` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::CognitoUserPoolUICustomization"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

These are the supported properties for the resource.

| Property name | Type | Description | Required |
| - | - | - | - |
| Domain | String | Domain name or part of the domain name to use | Yes |
| UserPoolId | String | The ID of the UserPool to create the Identity Provider in | Yes |
| CustomDomainConfig | Object | Object with CustomDomainConfig | Yes if custom domain is used |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

See more on [https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-ui-customization.html](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-ui-customization.html)

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
    Type: "Custom::CognitoUserPoolClient"
    DependsOn:
      - "UserPool"
    Properties:
      ClientName: "testclient"
      SupportedIdentityProviders:
        - "COGNITO"
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-client-${AWS::Region}-${Environment}"
      UserPoolId: !Ref "UserPool"

  CognitoClientUI:
    Type: "Custom::CognitoUserPoolUICustomization"
    DependsOn:
      - "UserPoolClient"
    Properties:
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-uicustomization-${AWS::Region}-${Environment}"
      CSS: ".logo-customizable {max-width: 100%; max-height: 40%;}"
      ImageFile: "iVBORw0KGgoAAAANSUhEUgAAAMgAAACnCAYAAABU+hMRA....=="
      ClientId: !Ref "UserPool"
      UserPoolId: !GetAtt "UserPoolClient.ClientId"
```