# custom-cf

Contains custom resources for creating resources in CloudFormation that currently doesn't exists.  
For creating your own please have a look in the `lib/` for common interfaces that can make your
life easier.

All these resources support all the cli attributes, but in CloudFormation.  
You can also use `Fn::GetAtt` on most of these to get data back from the resource.  
Please look at the individual `README.md` files in the functions folders.

## Custom Resources

### Cognito

- [cognito/userpool-federation](cognito/userpool-federation)
- [cognito/userpool-uicustomization](cognito/userpool-uicustomization)
- [cognito/userpool-client](cognito/userpool-client)

### ECS

- [ecs/tags](ecs/tags)

### IAM
- [iam/role-tags](iam/role-tags)


## Requirements

- docker
- aws cli

## Deployment

Use the included `Makefile` to deploy the resources.

The `OWNER` env var is for tagging. So you can set this to what you want.

```bash
AWS_PROFILE=my-profile AWS_REGION=region OWNER=TeamName S3_BUCKET=my-artifact-bucket FUNCTION=folder/my-resource make deploy
```

Example

```bash
AWS_PROFILE=default AWS_REGION=eu-west-1 OWNER=devops S3_BUCKET=my-artifact-bucket FUNCTION=cognito/userpool-federation deploy
```

## Creating a new Custom Resource

To create a new custom resource, please have a look in `example` folder for a simple example custom resource.
It will use the `lib/events` package (please see info about it below).
