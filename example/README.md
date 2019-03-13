# example

Basic boiletplate code for creating a custom resource for CloudFormation.

## Folder structure / naming

The function should be groupd by major AWS service, such as `cognito`, `s3` etc. And then in subfolders with a descrptive
name that best describes the API calls.

such as: `cognito/userpool-federation`.

## What needs to be changed

### template.yaml

In the `template.yaml` file you only need to change `Description` as well as the `Policy` document of the Role.  
All other values can be kept as is and will be automatically set by the `Makefile`.


### source code / main.go

You will need to change some of the constants at the top of the file.  
`service` can be left as is, if you haven't changed the whole project name in the `Makefile`.

`httpTimeout` can also be left as is if you don't have any special requirements with API calls that need longer
time to finish or needs a more fast timeout.

```go
const (
    service      = "custom-cf"
    function     = "myresource"         // Replace with the function name
    resourceType = "Custom::MyResource" // Change to the Resource Name you want to use.
    httpTimeout  = 30
)
```

`svc` in `config` struct should be replaced with an AWS Service (if used).

```go
    log *l.Client
    svc interface{} // Replace with AWS service (or other service etc) that the resource needs access to.

    physicalID            string              // The physical ID to use for the resource.
    resourceProperties    *ResourceProperties // The new resource data from the template.
    oldResourceProperties *ResourceProperties // The old resource data, only on updates.
}
```

Then the `ResourceProperties` struct needs to be updated to correspond to the data you have in your Custom Resource.

```go
// ResourceProperties needs to be exported so that the lib/events package can Unmarshal it.
// This should contain all the fields that you can add in the Custom resource in CloudFormation.
type ResourceProperties struct {
    MyResourceField1 string `json:"MyResourceField1"`
    MyResourceField2 string `json:"MyResourceField2"`
}
```

In the `handler` function the following lines should be replaced with creating the AWS (or other) service and
setting it to the c.svc field.

```go
// Function to create the AWS service (if needed) and set it to c.svc.
// If creation of the service fails it should return c.runError(req, err).
```

In the `run` function the following lines should be replaced with logic/function for determining if the resource
requested already exsits or not.

```go
// Add logic for checking if the resource with the same data already exists.
// This is just a placeholder variable.
exists := false
```

Then after that, the only thing that needs to be done is to add the logic/functions for the different steps like
create, update, delete (and depending on if the resource exists or not).
The map[string]string that is returned is only for create / update and is key values that can be accessed in
CloudFormation by the `Fn::GetAtt` function.

```go
switch {
// If Delete is run on the stack but the resource doesn't exist / already deleted.
case req.RequestType == "Delete" && !exists:
    return nil, nil

// If Delete is run on the stack.
case req.RequestType == "Delete" && exists:
    // Add logic to delete resource here.
    err := fmt.Errorf("placeholder result")
    return nil, err

// If Update is run on the stack but the resource doesn't exists
// create it. If it was a resource that needed replacement a delete event
// will be sent on the old resource once the new one has been created.
case req.RequestType == "Update" && !exists:
    // Add logic to create resource here.
    err := fmt.Errorf("placeholder result")
    return map[string]string{"key1": "value1"}, err

// If Update is run on the stack.
case req.RequestType == "Update" && exists:
    // Add logic to update resource here.
    err := fmt.Errorf("placeholder result")
    return map[string]string{"key1": "value1"}, err

// If Create is run on the stack and the resource doesn't exist:
case req.RequestType == "Create" && !exists:
    // Add logic to create resource here.
    err := fmt.Errorf("placeholder result")
    return map[string]string{"key1": "value1"}, err

// If Create is run on the stack and the resource exists, adopt and update it.
case req.RequestType == "Create" && exists:
    // Add logic to update resource here.
    err := fmt.Errorf("placeholder result")
    return map[string]string{"key1": "value1"}, err
}
```

## Deploying

To deploy simple use the included `Makefile` and run make deploy.  
For example

```bash
AWS_PROFILE=my-profile AWS_REGION=eu-west-1 OWNER=TeamName S3_BUCKET=my-artifact-bucket FUNCTION=cognito/my-resource make deploy
```