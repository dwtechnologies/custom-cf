# tags
Custom Cloudformation resource to support:
- AWS::ECS::Cluster.Tags
- AWS::ECS::TaskDefinition.Tags
- AWS::ECS::Service.Tags


Sample Cloudformation resource
```yaml
Resources:
  ECScluster:
    Type: AWS::ECS::Cluster
    Properties:
      ClusterName: production-cluster

  DefaultClusterTags:
    Type: Custom::ECSTags
    Properties:
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:<function-name>
      ResourceArn: !GetAtt ECScluster.Arn
      Tags:
       - Key: Location
         Value: stockholm
       - Key: Environment
         Value: prod
       - Key: Owner
         Value: cloudops

  AlpineTaskDefinitionTags:
    Type: Custom::ECSTags
    Properties:
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:<function-name>
      ResourceArn: !Sub arn:aws:ecs:${AWS::Region}:${AWS::AccountId}:task-definition/alpine:1
      Tags:
       - Key: Location
         Value: stockholm
       - Key: Environment
         Value: prod
       - Key: Owner
         Value: cloudops

  WebServiceTags:
    Type: Custom::ECSTags
    Properties:
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:<function-name>
      ResourceArn: !Sub arn:aws:ecs:${AWS::Region}:${AWS::AccountId}:service/default/web
      Tags:
       - Key: Location
         Value: stockholm
       - Key: Environment
         Value: prod
       - Key: Owner
         Value: cloudops
```

