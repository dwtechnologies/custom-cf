# userpool-client

Adds the ability to create, update and delete UserPool Client Settings through CloudFormation.

## Resource

The name for this custom resource is `Custom::CognitoUserPoolClient` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::CognitoUserPoolClient"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

These are the supported properties for the resource.

| Property name | Type | Description | Required |
| - | - | - | - |
| ClientName | String | The name of the Client. This is required by this implementation (but not in regular API!) | Yes |
| UserPoolId | String | The ID of the UserPool to create the Identity Provider in | Yes |
| GenerateSecret | bool | If we should generate secret. If you adopt a resource, make sure this setting is correct. Since changing this requires replacement on the client. Defaults to false. | No |
| RefreshTokenValidity | int | Token refresh validity | No |
| ReadAttributes | List of strings | Read Attributes | No |
| WriteAttributes | List of strings | Write Attributes | No |
| ExplicitAuthFlows | List of strings | Explicit Auth Flows | No |
| AllowedOAuthFlows | List of strings | Allowed OAuth Flows | No |
| AllowedOAuthFlowsUserPoolClient | String | Allowed OAuth Flows UserPool Client | No |
| AllowedOAuthScopes | List of strings | Allowed OAuth Scopes | No |
| CallbackURLs | List of strings | Callback URLs | No |
| LogoutURLs | List of strings | Logout URLs | No |
| DefaultRedirectURI | String | Default Redirect URI | No |
| SupportedIdentityProviders | List of strings | Name of supported providers (ProviderName). For current UserPool add `COGNITO`. | No |
| AnalyticsConfiguration | AnalyticsConfiguration | Analytics Configuration | No |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

For more details about userpool client check [https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-user-pool-client.html](https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/create-user-pool-client.html).

### AnalyticsConfiguration Properties

| Property name | Type | Description | Required |
| - | - | - | - |
| ApplicationId | String | Application Id | No |
| ExternalId | String | External ID | No |
| RoleArn | String | ARN to Role | No |
| UserDataShared | Bool | User Data Shared | No |

## Supported Attributes

The following attributes can be used in CloudFormations `Fn::GetAtt` function.

- ClientName
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
```