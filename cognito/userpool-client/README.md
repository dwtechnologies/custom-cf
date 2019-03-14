# userpool-client

Adds the ability to create, update and delete UserPool Client Settings through CloudFormation.

## Resource name

The name for this custom resource is `Custom::CognitoUserPoolClient` and
supports all the parameters that you can make through the GUI and cli.

For parameter list check all the parameters of the [https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-user-pool-client.html](https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-user-pool-client.html).
The parameters should be named the same as in the cli but with CamelCase instead of lower case and with hyphens.

If you have renamed the cognito-userpool-client function, please update the `ServiceToken` below to the corresponding lambda.

Please note that if you adopt an already existing client it's paramount that you choose the correct `GenerateSecret` setting as when
when it was created manually. Otherwise the resource might get deleted when you update it in the future due to requiring replacement.

## Supported Attributes

The following attributes can be used in CloudFormations `Fn::GetAtt` function.

- ClientName
- ClientId
- UserPoolId
- ClientSecret

## Example

```yaml
AWSTemplateFormatVersion: "2010-09-09"
Description: "Cognito UserPool with Client Settings"

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
      - "UserPoolFederationADFS"
    Properties:
      ClientName: "testclient"
      SupportedIdentityProviders:
        - "MyProvider"
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-client-${AWS::Region}-${Environment}"
      UserPoolId: !Ref "UserPool"
```