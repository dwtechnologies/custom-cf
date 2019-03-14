# custom-cf

Contains custom resources for creating resources in CloudFormation that currently doesn't exists.  
For creating your own please have a look in the `lib/` for common interfaces that can make your
life easier.

All these resources support all the cli attributes, but in CloudFormation.  
You can also use `Fn::GetAtt` on most of these to get data back from the resource.  
Please look at the individual `README.md` files in the functions folders.

## Custom Resources

- [cognito/userpool-federation](cognito/userpool-federation)
- [cognito/userpool-uicustomization](cognito/userpool-uicustomization)
- [cognito/userpool-client](cognito/userpool-client)

## Creating a new Custom Resource

To create a new custom resource, please have a look in `example` folder for a simple example custom resource.
It will use the `lib/events` package (please see info about it below).
