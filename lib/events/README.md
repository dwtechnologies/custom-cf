# events

Is used for parsing incoming Request to the lambda as well as creating the Response and sending it
to the s3 pre-signed URL.

You will need to call the Unmarshal method and supply a struct with JSON tags as new (ResourceProperties)
and old (OldResourcePropeties) to unmarshal the Custom properties. Since the data structure will differ depending on the
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
    "github.com/dwtechnologies/custom-cf/lib/events"
    "github.com/aws/aws-lambda-go/lambda"
)

type Resource struct {
    Key1 string `json:"Key1"`
    Key2 string `json:"Key2"`
}

func main() {
    lambda.Start(handler)
}

func handler(ctx context.Context, req *events.Request) error {
    // Unmarshal ResourceProperties and OldResourceProperties.
    new, old := &Resource{}, &Resource{}
    if err := req.Unmarshal(new, old); err != nil {
        if err := req.Send("testID1", nil, err); err != nil {
            return err
        }
        return err
    }

    // Do something here with the data...

    // Save the response to s3.
    return req.Send("testID1", map[string]string{"key1": "value1"}, nil)
}
```