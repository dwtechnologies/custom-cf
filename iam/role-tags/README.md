# tags
Custom Cloudformation resource to support:
- AWS::IAM::Role.Tags


Sample Cloudformation resource
```yaml
Resources:
  SomeRoleTags:
    Type: Custom::IAMRoleTags
    Properties:
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:<function-name>
      RoleName: somerole-name
      Tags:
       - Key: Location
         Value: stockholm
       - Key: Environment
         Value: prod
       - Key: Owner
         Value: cloudops
```

