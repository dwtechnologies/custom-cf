# custom-cf

Contains custom resources for creating resources in CloudFormation that currently doesn't exists.  
For creating your own please have a look in the `lib/` for common interfaces that can make your
life easier.

## libs

### lib/respond

Is used for parsing incoming Request to the lambda as well as creating the Response and sending it
to the s3 pre-signed URL.

The `ResourceProperties` and `OldResourceProperties` will be in RAW JSON format. So you will need to
manually Unmarshal them into your own structs. Since the data structure will differ depending on the
resource types you want to create a Custom Resource for.

Example usage below.
We will save the resource with a "physicalId" (should be unique) `testID1` and with data that can be
accessed through the `Fn::GetAtt` function for the key value `key1` with value `value1`.

Error will be nil, since we didn't get an error when creating the resource in our lambda function.  
Here you would send the resource creation error so that the error state and reason why it failed
will be saved in CF and can be visible in the console.

```go
package main

import (
    "context"
    "github.com/dwtechnologies/custom-cf/lib/respond"
    "github.com/aws/aws-lambda-go/lambda"
)

func main() {
    lambda.Start(handler)
}

func handler(ctx context.Context, req *respond.Request) error {
    // Do something with req.
    // Unmarshal ResourceProperties and OldResourceProperties.
    return req.Send("testID1", map[string]string{"key1": "value1"}, nil)
}
```