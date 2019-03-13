# cognitoidentityprovider/ui-customization


```yaml
Resources:
  CognitoClientUIWengy:
    Type: Custom::CognitoUserPoolUICustomization
    Properties:
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::Account}:function:<function-name>
      CSS: ".logo-customizable {\n\tmax-width: 100%;\n\tmax-height: 40%;\n}"
      ClientId: 672ks333jjkll3
      UserPoolId: eu-west-1_Gxf890

```

