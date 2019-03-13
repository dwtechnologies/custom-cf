# cognito/userpool-uicustomization

Sample cloudformation resource
```yaml
Resources:
  CognitoClientUI:
    Type: Custom::CognitoUserPoolUICustomization
    Properties:
      ServiceToken: !Sub arn:aws:lambda:${AWS::Region}:${AWS::Account}:function:<function-name>
      CSS: ".logo-customizable {max-width: 100%; max-height: 40%;}"
      ImageFile: "base64 encoded logo. up to 100 KB in size."
      ClientId: 672ks333jjkll3
      UserPoolId: eu-west-1_Gxf890
```

#### CSS customization reference:
https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-ui-customization.html

