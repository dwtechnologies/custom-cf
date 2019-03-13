# userpool-federation

Adds the ability to attach UserPool Federation through CloudFormation.

## Resource name

The name for this custom resource is `Custom::CognitoUserPoolFederation` and
supports all the parameters that you can make through the GUI and cli.

For parameter list check all the parameters of the [https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-identity-provider.html](https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-identity-provider.html).
The parameters should be named the same as in the cli but with CamelCase instead of lower case and with hyphens.

If you have renamed the cognito-userpool-federation function, please update the `ServiceToken` below to the corresponding lambda.

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
      ProviderName: "Adopt"
      ProviderType: "SAML"
      ProviderDetails:
        MetadataURL: "https://my.domain.com/FederationMetadata.xml"
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-userpool-federation-${AWS::Region}-${Environment}"
      UserPoolId: !Ref "UserPool"
```