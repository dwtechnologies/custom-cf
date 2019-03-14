# ecs/tags

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
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::Account}:function:<function-name>
      ResourceARN: !GetAtt ECScluster.Arn
      Tags:
	- Key: Location
	  Value: Stockholm
	- Key: Environment
	  Value: prod
	- Key: Owner
	  Value: cloudops
```

