# ecs/tags
Custom cloudformation resource to support
- AWS::ECS::Cluster.Tags
- AWS::ECS::TaskDefinition.Tags
- AWS::ECS::Service.Tags


Sample cloudformation resource
```yaml
Resources:
  ECScluster:
    Type: AWS::ECS::Cluster
    Properties:
      ClusterName: production-cluster

  DefaultClusterTags:
    Type: Custom::ECSTag
    Properties:
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:<function-name>
      ResourceArn: !GetAtt ECScluster.Arn
      Tags:
	- Key: Location
	  Value: Stockholm
	- Key: Environment
	  Value: prod
	- Key: Owner
	  Value: cloudops
```

