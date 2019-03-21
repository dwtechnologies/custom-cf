# tags

Adds the ability to set tags for `ECS::Cluster`, `ECS::TaskDefinition` and `ECS::Service` through CloudFormation.

## Resource

The name for this custom resource is `Custom::ECSTags` and
supports all the parameters that you can make through the GUI and cli.

## Structure

This is the YAML structure you use when using this Custom Resource.

```yaml
Type: "Custom::ECSTags"
Properties:
  Properties
```

See below for the supported Properties.

## Properties

These are the supported properties for the resource.

| Property name | Type | Description | Required |
| - | - | - | - |
| ResourceArn | String | ARN to the ECS Cluster, Task Definition or Service | Yes |
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
  ECScluster:
    Type: "AWS::ECS::Cluster"
    Properties:
      ClusterName: "production-cluster"

  DefaultClusterTags:
    Type: "Custom::ECSTags"
    Properties:
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-ecs-tags-${AWS::Region}-${Environment}"
      ResourceArn: !GetAtt "ECScluster.Arn"
      Tags:
       - Key: "Location"
         Value: "stockholm"
       - Key: "Environment"
         Value: "prod"
       - Key: "Owner"
         Value: "cloudops"

  AlpineTaskDefinitionTags:
    Type: "Custom::ECSTags"
    Properties:
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-ecs-tags-${AWS::Region}-${Environment}"
      ResourceArn: !Sub "arn:aws:ecs:${AWS::Region}:${AWS::AccountId}:task-definition/alpine:1"
      Tags:
       - Key: "Location"
         Value: "stockholm"
       - Key: "Environment"
         Value: "prod"
       - Key: "Owner"
         Value: "cloudops"

  WebServiceTags:
    Type: "Custom::ECSTags"
    Properties:
      ServiceToken: !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:cognito-ecs-tags-${AWS::Region}-${Environment}"
      ResourceArn: !Sub "arn:aws:ecs:${AWS::Region}:${AWS::AccountId}:service/default/web"
      Tags:
       - Key: "Location"
         Value: "stockholm"
       - Key: "Environment"
         Value: "prod"
       - Key: "Owner"
         Value: "cloudops"
```
