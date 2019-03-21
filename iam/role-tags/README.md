# tags

Adds the ability to set tags for `IAM::Role` through CloudFormation.

## Resource

The name for this custom resource is `Custom::IAMRoleTags` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::IAMRoleTags"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

These are the supported properties for the resource.

| Property name | Type | Description | Required |
| - | - | - | - |
| RoleName | String | Role name | Yes |
| Tags | Tags | List of tags | Yes |
| ServiceToken | String | The ARN of the lambda function for this Custom Resource | Yes |

### Tag properties

| Property name | Type | Description | Required |
| - | - | - | - |
| Key | String | Name of tag | Yes |
| Value | String | Value of tag | Yes |

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
  IamRoleTags:
    Type: "Custom::ECSTags"
    Properties:
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-iam-tags-${AWS::Region}-${Environment}"
      RoleName: "somerole-name"
      Tags:
       - Key: "Location"
         Value: "stockholm"
       - Key: "Environment"
         Value: "prod"
       - Key: "Owner"
         Value: "cloudops"
```
