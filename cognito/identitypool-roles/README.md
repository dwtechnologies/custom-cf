# identitypool-roles

Adds the ability to set Identity Pool Role mappings based on dynamic Identity Provider names.  
Since this is not possible in the included Identity Pool Role Attachment of CloudFormation.

Which makes it impossible to do Role Attachment on UserPool Clients created in CloudFormation.

This function will mitigate that limitation.

## Resource

The name for this custom resource is `Custom::CognitoIdentityPoolRoles` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::CognitoIdentityPoolRoles"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

| Property name | Type | Description | Required |
| - | - | - | - |
| Roles | Map of strings | Set the default roles ARNs. Currently supports keys **authenticated** and **unauthenticated** | No |
| RoleMappings | List of RoleMapping | The role mapping for a specific Identity Provider | No |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

For more details about the properties check the aws cli docs [https://docs.aws.amazon.com/cli/latest/reference/cognito-identity/set-identity-pool-roles.html](https://docs.aws.amazon.com/cli/latest/reference/cognito-identity/set-identity-pool-roles.html).

### RoleMapping Properties

When using UserPool the Identity Provider should be `cognito-idp.${Region}.amazonaws.com/${UserPoolId}:${UserPoolClientId}`.

Up to 25 rules can be specified per Identity Provider.

| Property name | Type | Description | Required |
| - | - | - | - |
| IdentityProvider | String | The Identity Provider to use for the mapping | Yes |
| Type | String | Can be either **Token** or **Rules**. Where **Token** will use `cognito:roles` and `cognito:preferred_role` claims. Where `Rules` will match the claims from the token to a role | Yes |
| AmbiguousRoleResolution | String | Specify the action to be taken if there is no match if using type **Rules** or there is no preferred_role when using **Token** and multiple roles. Valid values are **AuthenticatedRole** or **Deny** | Yes |
| Rules | List of Rule | List of rules. This is required when choosing type **Rules** | No |

### Rule Properties

The rules for mapping roles if type is **Rules**.

| Property name | Type | Description | Required |
| - | - | - | - |
| Claim | String | The Claim to match | Yes |
| MatchType | String | Either **Equals**, **Contains**, **StarsWith** or **NotEqual** | No |
| Value | String | The value to match against the claim | Yes |
| RoleArn | String | The ARN to the role to assign | Yes |

## Example

The following example will add an default authenticated role called `AuthenticatedRole` and map against a role called `MappedRole`
if claim `testGroup` is `testGroup`.

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

  IdentityPool:
    Type: "AWS::Cognito::IdentityPool"
    DependsOn:
      - "UserPoolClient"
    Properties:
      IdentityPoolName: !Sub "identitypool"
      AllowUnauthenticatedIdentities: false
      CognitoIdentityProviders:
        - ClientId: !GetAtt "UserPoolClient.ClientId"
          ProviderName: !GetAtt "UserPool.ProviderName"
          ServerSideTokenCheck: false

  IdentityPoolRoleMappings:
    Type: "Custom::CognitoIdentityPoolRoles"
    DependsOn:
      - "IdentityPool"
      - "AuthenticatedRole"
      - "MappedRole"
    Properties:
      IdentityPoolId: !Ref "IdentityPool"
      Roles:
        authenticated: !GetAtt "AuthenticatedRole.Arn"
      RoleMappings:
        - IdentityProvider: !Sub "cognito-idp.${AWS::Region}.amazonaws.com/${UserPool}:${UserPoolClient.ClientId}"
          Type: "Rules"
          AmbiguousRoleResolution: "Deny"
          Rules:
            - Claim: "group"
              MatchType: "Equals"
              Value: "testGroup"
              RoleArn: !GetAtt "MappedRole.Arn"
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-identitypool-roles-${AWS::Region}-${Environment}"
```