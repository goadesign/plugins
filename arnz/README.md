# Arnz

ArnZ is a DSL for authorizing methods based on [AWS IAM](https://aws.amazon.com/iam/) caller ARNs.

## Given

Your Goa application...
1. is recieving traffic via an [AWS API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api.html).
1. is using the [AWS_IAM](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-access-control-iam.html) authorizer.

## You Can

### Authenticate All Callers

When imported, all methods will require all callers to be IAM authenticated.

```go
package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/arnz/dsl"
)
```

### Authorize Callers by ARN

You can authorize callers by ARN using the `AllowArnsMatching` function, passing it a regular expression. 

```go
Method("privileged", func() {
	AllowArnsMatching("^arn:aws:iam::123456789012:user/administrator$")
	Result(SecretStuff)
	HTTP(func() {
		Get("/secrets")
		Response(StatusOK)
	})
})
```

### Allow Unsigned Requests

Allowing unsigned requests is useful for allowing traffic not originated from API gateway. 

```go
Method("healthz", func() {
	AllowUnsignedCallers()
	Result(HealthCheck)
	HTTP(func() {
		GET("/healthz")
		Response(StatusOK)
	})
})
```

_note_: Allowing unsigned callers does not disable authentication or authorization for signed requests.

## Further Reading
- [Signing HTTP requests using AWS credentials](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_aws-signing.html)
- [API Gateway Developer Docs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api.html)
- [API Gateway IAM Authorizer Docs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-access-control-iam.html)